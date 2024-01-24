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
	"strconv"
	"strings"
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

	files := getSortedXmlFiles(zf)
	filesWithAmount := shutdown.GetFilesWithAmount()

	threadNumber := 3
	if tn := config.GetConfig("APP_THREAD_NUMBER"); tn != "" {
		threadNumber, _ = strconv.Atoi(tn)
	}
	mutexChan := make(chan struct{}, threadNumber)
	g, onErrCtx := errgroup.WithContext(context.Background())
	for _, file := range files {
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
		f, err := _file.Open()
		g.Go(func() error {
			select {
			case <-onErrCtx.Done():
				<-mutexChan
				return nil
			default:
				var amountForDump int
				fileName := _file.Name
				amount, err := ProcessingXml(
					f,
					objectType,
					func(ol *types.FiasObjectList) error {
						amountInFile, ok := filesWithAmount[fileName]
						var amountForDumpResult int
						if ok && amountInFile > amountForDump {
							amountForDumpResult = amountInFile
						} else {
							amountForDumpResult = amountForDump
						}
						select {
						case <-onErrCtx.Done():
							shutdown.PutFileToDump(shutdown.DumpFile{FileName: fileName, RecordsAmount: amountForDumpResult})
							return errors.New("error when import, thread stop")
						case <-ctx.Done():
							shutdown.PutFileToDump(shutdown.DumpFile{FileName: fileName, RecordsAmount: amountForDumpResult})
							return errors.New("shutdown, thread stop")
						default:
							if ok && amountInFile > amountForDump {
								amountForDump += len(ol.List)
								return nil
							}
							err = importToDb(ol)
							if err == nil {
								amountForDump += len(ol.List)
							} else {
								shutdown.PutFileToDump(shutdown.DumpFile{FileName: fileName, RecordsAmount: amountForDumpResult})
							}
							return err
						}
					},
				)
				if err != nil {
					<-mutexChan
					return err
				}
				<-mutexChan
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
