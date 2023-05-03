package fias

import (
	"encoding/xml"
	"fias_to_sql/internal/services/fias/types"
	"io"
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
				}
				err := decoder.DecodeElement(&fiasObj, &se)
				if err != nil {
					return nil, err
				}
				if !fiasObj.GetIsActive() {
					continue
				}
				al.AddObject(fiasObj)
				//Todo debug
				return al, nil
			}
		}
	}
	return al, nil
}
