package fias

import (
	"encoding/xml"
	"fmt"
	"io"
)

type fiasObject interface {
	getFiasObjectId() int64
}

type address struct {
	fiasObject
	Id         int64  `xml:"ID,attr"`
	ObjectId   int64  `xml:"OBJECTID,attr"`
	ObjectGuid string `xml:"OBJECTGUID,attr"`
	TypeName   string `xml:"TYPENAME,attr"`
	Level      int64  `xml:"LEVEL,attr"`
	Name       string `xml:"NAME,attr"`
}

type house struct {
	fiasObject
	Id       int64 `xml:"ID,attr"`
	ObjectId int64 `xml:"OBJECTID,attr"`
	HouseNum int64 `xml:"HOUSENUM,attr"`
	AddNum   int64 `xml:"ADDNUM,attr"`
	AddNum2  int64 `xml:"ADDNUM2,attr"`
}

type hierarchy struct {
	fiasObject
}

type FiasObjectsList struct {
	addresses []fiasObject
}

func (a FiasObjectsList) addObject(object fiasObject) {
	a.addresses = append(a.addresses, object)
}

func (a FiasObjectsList) clear() {
	a.addresses = make([]fiasObject, 0)
}

func ProcessingXml(
	closer io.ReadCloser,
	object ...string,
) (*FiasObjectsList, error) {
	defer closer.Close()
	objectType := "OBJECT"
	if len(object) > 0 {
		objectType = object[0]
	}
	decoder := xml.NewDecoder(closer)
	al := new(FiasObjectsList)
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
			if se.Name.Local == objectType {
				var fiasObj fiasObject
				switch objectType {
				case "OBJECT":
					fiasObj = &address{}
				}
				err := decoder.DecodeElement(&fiasObj, &se)
				if err != nil {
					return nil, err
				}
				fmt.Println(fiasObj)
			}
		}
	}
	return al, nil
}
