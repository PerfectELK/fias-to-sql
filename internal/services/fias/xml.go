package fias

import (
	"encoding/xml"
	"io"
)

type address struct {
	Id       int64
	ParentId int64
	Name     string
}

type AddressList struct {
	addresses []address
}

func (a AddressList) addAddress(addr address) {
	a.addresses = append(a.addresses, addr)
}

func ProcessingXml(closer io.ReadCloser) (*AddressList, error) {
	defer closer.Close()
	d := xml.NewDecoder(closer)
	al := new(AddressList)
	for {
		t, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			return nil, tokenErr
		}
	}
	return al, nil
}
