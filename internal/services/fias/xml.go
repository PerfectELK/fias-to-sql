package fias

import (
	"encoding/xml"
	"fias_to_sql/internal/services/fias/types"
	"io"
	"time"
)

func ProcessingXml(
	closer io.ReadCloser,
	objectType string,
	fn func(ol *types.FiasObjectList) error,
) (int, error) {
	defer closer.Close()

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
	case "obj-types":
		xmlTag = "ADDRESSOBJECTTYPE"
	}

	decoder := xml.NewDecoder(closer)
	al := new(types.FiasObjectList)
	var counter int
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			return counter, tokenErr
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
				case "obj-types":
					fiasObj = &types.AddressObjectType{}
				}

				err := decoder.DecodeElement(&fiasObj, &se)
				if err != nil {
					return 0, err
				}

				switch fo := fiasObj.(type) {
				case *types.Address, *types.House, *types.Hierarchy:
					fieldProcessing(al, &fiasObj)
				case *types.Param:
					paramProcessing(al, fo)
				case *types.AddressObjectType:
					objTypeProcessing(al, &fiasObj)
				}

				if len(al.List) >= 2000 {
					err := fn(al)
					if err != nil {
						return 0, err
					}
					counter += len(al.List)
					al.Clear()
				}
			}
		}
	}
	err := fn(al)
	if err != nil {
		return 0, err
	}
	counter += len(al.List)
	al.Clear()

	return counter, nil
}

func ProcessingXmlToChan(
	closer io.ReadCloser,
	objectType string,
	ch chan<- *types.FiasObjectList,
	passAmount ...int,
) (int, error) {
	pass := 0
	if len(passAmount) > 0 {
		pass = passAmount[0]
	}
	defer closer.Close()

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
	case "obj-types":
		xmlTag = "ADDRESSOBJECTTYPE"
	}

	decoder := xml.NewDecoder(closer)
	al := new(types.FiasObjectList)
	var counter int
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			return counter, tokenErr
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
				case "obj-types":
					fiasObj = &types.AddressObjectType{}
				}

				err := decoder.DecodeElement(&fiasObj, &se)
				if err != nil {
					return 0, err
				}

				switch fo := fiasObj.(type) {
				case *types.Address, *types.House, *types.Hierarchy:
					fieldProcessing(al, &fiasObj)
				case *types.Param:
					paramProcessing(al, fo)
				case *types.AddressObjectType:
					objTypeProcessing(al, &fiasObj)
				}

				if pass > 0 {
					pass -= len(al.List)
					al.Clear()
				}

				if len(al.List) >= 2000 {
					ch <- al
					counter += len(al.List)
					al = new(types.FiasObjectList)
				}
			}
		}
	}
	ch <- al
	counter += len(al.List)
	close(ch)

	return counter, nil
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

func objTypeProcessing(
	list *types.FiasObjectList,
	field *types.FiasObject,
) {
	list.AddObject(*field)
}
