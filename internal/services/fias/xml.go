package fias

import (
	"encoding/xml"
	"fias_to_sql/internal/services/fias/types"
	"io"
	"time"
)

func ProcessingXml(
	closer io.ReadCloser,
	object ...string,
) (*types.FiasObjectList, error) {
	defer closer.Close()
	objectType := "object"
	if len(object) > 0 {
		objectType = object[0]
	}

	var xmlTag string
	switch objectType {
	case "object":
		xmlTag = "OBJECT"
	case "house":
		xmlTag = "HOUSE"
	case "hierarchy":
		xmlTag = "ITEM"
	case "param":
		xmlTag = "PARAM"
	}

	decoder := xml.NewDecoder(closer)
	al := new(types.FiasObjectList)
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			return nil, tokenErr
		}
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == xmlTag {
				var fiasObj types.FiasObject
				switch objectType {
				case "object":
					fiasObj = &types.Address{}
				case "house":
					fiasObj = &types.House{}
				case "hierarchy":
					fiasObj = &types.Hierarchy{}
				case "param":
					fiasObj = &types.Param{}
				}

				err := decoder.DecodeElement(&fiasObj, &se)
				if err != nil {
					return nil, err
				}

				switch fo := fiasObj.(type) {
				case *types.Address, *types.House, *types.Hierarchy:
					fieldProcessing(al, &fiasObj)
				case *types.Param:
					paramProcessing(al, fo)
				}
			}
		}
	}
	return al, nil
}

func fieldProcessing(
	list *types.FiasObjectList,
	field *types.FiasObject,
) {
	if !(*field).GetIsActive() {
		return
	}
	list.AddObject(*field)
}

func paramProcessing(
	list *types.FiasObjectList,
	param *types.Param,
) {
	if param.TypeId != 10 {
		return
	}

	endTime, err := time.Parse("2006-01-02", param.EndDate)
	if err != nil {
		return
	}
	now := time.Now()
	if now.Unix() > endTime.Unix() {
		return
	}
	startTime, err := time.Parse("2006-01-02", param.StartDate)
	if err != nil {
		return
	}
	if now.Unix() < startTime.Unix() {
		return
	}

	list.AddObject(param)
}
