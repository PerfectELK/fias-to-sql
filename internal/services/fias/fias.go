package fias

import (
	"archive/zip"
	"context"
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/models"
	"fias_to_sql/internal/services/fias/types"
	"fias_to_sql/internal/services/logger"
	"fias_to_sql/internal/services/shutdown"
	"fias_to_sql/pkg/slice"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
	"sync/atomic"
)

func getSortedXmlFiles(zf *zip.ReadCloser) []*zip.File {
	zipFiles := make([]*zip.File, 0)
	shutdownFiles := shutdown.GetFilesNames()
	for _, file := range zf.File {
		if shutdown.IsReboot && !slice.Contains(shutdownFiles, file.Name) {
			continue
		}

		var objectType string
		if strings.Contains(file.Name, config.GetConfig("HOUSES_FILE_PART")) {
			objectType = "house"
			if strings.Contains(file.Name, "AS_HOUSES_PARAMS") {
				continue
			}
		}
		if strings.Contains(file.Name, config.GetConfig("OBJECT_FILE_PART")) {
			objectType = "object"
		}
		if strings.Contains(file.Name, config.GetConfig("HIERARCHY_FILE_PART")) {
			objectType = "hierarchy"
		}
		if strings.Contains(file.Name, "_OBJ_TYPES_") {
			objectType = "obj-types"
		}
		if strings.Contains(file.Name, "_OBJ_PARAMS_") {
			objectType = "param"
		}
		if objectType == "" {
			continue
		}
		if strings.Contains(file.Name, "_DIVISION_") {
			continue
		}

		zipFiles = append(zipFiles, file)
	}
	return zipFiles
}

var filesWithAmount map[string]int

func ImportXml(
	ctx context.Context,
	archivePath string,
) error {
	shutdown.SetArchivePathToDump(archivePath)

	zf, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zf.Close()

	files := filesToChan(getSortedXmlFiles(zf))

	filesWithAmount = shutdown.GetFilesWithAmount()

	g, onErrCtx := errgroup.WithContext(context.Background())
	mutexChan := make(chan struct{}, 20)
	for file := range files {
		if ctx.Err() != nil {
			shutdown.PutFileToDump(shutdown.DumpFile{FileName: file.Name, RecordsAmount: 0})
			continue
		}

		var objectType string
		if strings.Contains(file.Name, config.GetConfig("HOUSES_FILE_PART")) {
			objectType = "house"
		}
		if strings.Contains(file.Name, config.GetConfig("OBJECT_FILE_PART")) {
			objectType = "object"
		}
		if strings.Contains(file.Name, config.GetConfig("HIERARCHY_FILE_PART")) {
			objectType = "hierarchy"
		}
		if strings.Contains(file.Name, "_PARAMS_") {
			objectType = "param"
		}
		if strings.Contains(file.Name, "_OBJ_TYPES_") {
			objectType = "obj-types"
		}

		mutexChan <- struct{}{}
		_file := file
		f, _ := _file.Open()
		g.Go(func() error {
			defer func() {
				<-mutexChan
			}()
			select {
			case <-onErrCtx.Done():
				return errors.New("error when import, thread stop")
			case <-ctx.Done():
				return errors.New("shutdown, thread stop")
			default:
				amountInFile, _ := filesWithAmount[_file.Name]
				fiasCh := make(chan *types.FiasObjectList, 100)

				var amountForDumpResult int64
				wg := sync.WaitGroup{}
				for i := 0; i < 30; i++ {
					wg.Add(1)
					g.Go(func() error {
						defer wg.Done()
						select {
						case <-onErrCtx.Done():
							shutdown.PutFileToDump(shutdown.DumpFile{FileName: _file.Name, RecordsAmount: int(amountForDumpResult)})
							return errors.New("error when import, thread stop")
						case <-ctx.Done():
							shutdown.PutFileToDump(shutdown.DumpFile{FileName: _file.Name, RecordsAmount: int(amountForDumpResult)})
							return errors.New("shutdown, thread stop")
						default:
							err := processingObjectList(fiasCh, &amountForDumpResult)
							if err != nil {
								shutdown.PutFileToDump(shutdown.DumpFile{FileName: _file.Name, RecordsAmount: int(amountForDumpResult)})
							}
							return err
						}
					})
				}

				amount, err := ProcessingXmlToChan(
					f,
					objectType,
					fiasCh,
					amountInFile,
				)

				wg.Wait()
				if err != nil {
					return err
				}
				logger.Println(_file.Name, ": records amount (", amount, ") [OK]")
				return nil
			}
		})
	}

	err = g.Wait()
	if err != nil {
		return err
	}

	return err
}

func processingObjectList(ch <-chan *types.FiasObjectList, counter *int64) error {
	for ol := range ch {
		err := importToDb(ol)
		if err != nil {
			return err
		}
		atomic.AddInt64(counter, int64(len(ol.List)))
		ol.Clear()
	}
	return nil
}

func filesToChan(zf []*zip.File) <-chan *zip.File {
	ch := make(chan *zip.File, 5)
	go func() {
		defer close(ch)
		for _, f := range zf {
			ch <- f
		}
	}()
	return ch
}

func importToDb(list *types.FiasObjectList) error {
	if len(list.List) == 0 {
		return nil
	}
	var modelList models.ModelListStruct
	modelList.List = make([]models.Model, 0, len(list.List))

	for _, item := range list.List {
		var err error
		switch fiasObj := item.(type) {
		case *types.Address:
			model := models.NewObject()
			model.SetName(fiasObj.Name)
			model.SetObject_id(fiasObj.ObjectId)
			model.SetObject_guid(fiasObj.ObjectGuid)
			model.SetLevel(fiasObj.Level)
			model.SetType_name(fiasObj.TypeName)
			modelList.AppendModel(model)
		case *types.House:
			model := models.NewObject()
			model.SetName(fiasObj.HouseNum)
			model.SetObject_id(fiasObj.ObjectId)
			model.SetObject_guid(fiasObj.ObjectGuid)
			model.SetLevel(12)
			model.SetType_name("дом")
			modelList.AppendModel(model)
		case *types.Hierarchy:
			model := models.NewHierarchy()
			model.SetObject_id(fiasObj.ObjectId)
			model.SetParent_object_id(fiasObj.ParentObjId)
			modelList.AppendModel(model)
		case *types.Param:
			model := models.NewParam()
			model.SetObject_id(fiasObj.ObjectId)
			model.SetKladr_id(fiasObj.Value)
			modelList.AppendModel(model)
		case *types.AddressObjectType:
			model := models.NewObjectType()
			model.SetId(fiasObj.Id)
			model.SetName(fiasObj.Name)
			model.SetShortName(fiasObj.ShortName)
			model.SetLevel(fiasObj.Level)
			modelList.AppendModel(model)
		}
		if err != nil {
			return err
		}
	}
	err := modelList.SaveModelList()
	if err != nil {
		return err
	}
	return nil
}
