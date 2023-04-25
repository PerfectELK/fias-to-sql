package types

type FiasObject interface {
	GetIsActive() bool
}

type Address struct {
	FiasObject
	ObjectId   int64  `xml:"OBJECTID,attr"`
	ObjectGuid string `xml:"OBJECTGUID,attr"`
	TypeName   string `xml:"TYPENAME,attr"`
	Level      int64  `xml:"LEVEL,attr"`
	Name       string `xml:"NAME,attr"`
	IsActive   bool   `xml:"ISACTIVE,attr"`
}

type House struct {
	FiasObject
	ObjectId   int64  `xml:"OBJECTID,attr"`
	ObjectGuid string `xml:"OBJECTGUID,attr"`
	HouseNum   string `xml:"HOUSENUM,attr"`
	AddNum     string `xml:"ADDNUM,attr"`
	AddNum2    string `xml:"ADDNUM2,attr"`
	IsActive   bool   `xml:"ISACTIVE,attr"`
}

type Hierarchy struct {
	FiasObject
	Id          int64 `xml:"ID,attr"`
	ObjectId    int64 `xml:"OBJECTID,attr"`
	ParentObjId int64 `xml:"PARENTOBJID,attr"`
	IsActive    bool  `xml:"ISACTIVE,attr"`
}

func (f Address) GetIsActive() bool {
	return f.IsActive
}
func (f House) GetIsActive() bool {
	return f.IsActive
}
func (f Hierarchy) GetIsActive() bool {
	return f.IsActive
}

type FiasObjectList struct {
	Addresses []FiasObject
}

func (a *FiasObjectList) AddObject(object FiasObject) {
	a.Addresses = append(a.Addresses, object)
}

func (a *FiasObjectList) Clear() {
	a.Addresses = nil
}
