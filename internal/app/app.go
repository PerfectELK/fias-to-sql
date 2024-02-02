package app

import (
	"context"
	"errors"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/internal/services/dirs"
	"github.com/PerfectELK/go-import-fias/internal/services/disk"
	"github.com/PerfectELK/go-import-fias/internal/services/error/handler"
	"github.com/PerfectELK/go-import-fias/internal/services/fias"
	"github.com/PerfectELK/go-import-fias/internal/services/logger"
	"github.com/PerfectELK/go-import-fias/internal/services/shutdown"
	"github.com/PerfectELK/go-import-fias/internal/services/terminal"
	"github.com/PerfectELK/go-import-fias/migrations"
	"github.com/PerfectELK/go-import-fias/pkg/db"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func Run() error {
	err := dirs.InitServiceDirs()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("begin init app")
	err = config.InitConfig()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	err = terminal.ParseArgs()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	usageGB, err := disk.Usage()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	if usageGB.FreeGB < 70 {
		return errors.New("no space left on device")
	}
	logger.Println("init app success")

	if shutdown.CheckGracefulShutdown() {
		logger.Println("reboot after graceful shutdown")
		err := shutdown.RebootAfterGracefulShutdown()
		if err != nil {
			return handler.ErrorHandler(err)
		}
	}

	path, err := fias.GetArchivePath()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	logger.Println("create db and tables if not exists")
	_, err = db.GetDbInstance()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	err = migrations.CreateDatabase()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	err = migrations.CreateTables()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("create db and tables success")

	logger.Println("begin import")
	beginTime := time.Now()
	defer func() {
		logger.Println("import time ", int(time.Since(beginTime).Seconds())/60, " minutes")
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown.OnShutdown(func() {
		cancel()
		logger.Println("start shutdown")
	})

	if config.GetConfig("APP_DEBUG") == "true" {
		logger.Println("debugger start")
		go func() {
			err := http.ListenAndServe("localhost:8585", nil)
			if err != nil {
				logger.Println("error when start debugger")
			}
		}()
	}

	logger.Println("start import")
	err = fias.ImportXml(
		ctx,
		path,
	)

	if ctx.Err() != nil {
		err := shutdown.MakeDump()
		if err != nil {
			return handler.ErrorHandler(err)
		}
		logger.Println("end shutdown")
		os.Exit(-1)
	}

	if err != nil {
		return handler.ErrorHandler(err)
	}

	logger.Println("begin create indexes")
	err = migrations.CreateIndexes()
	if err != nil {
		return err
	}
	logger.Println("indexes created success")

	logger.Println("begin migrate data from temp to original tables")
	err = migrations.MigrateDataFromTempTables()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("migration success")

	logger.Println("create additional views")
	err = migrations.CreateAdditionalViews()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("create additional views success")
	logger.Println("import success")
	return nil
}
