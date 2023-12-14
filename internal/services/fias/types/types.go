package types

type FiasObject interface {
	GetIsActive() bool
}

type Address struct {
	ObjectId   int64  `xml:"OBJECTID,attr"`
	ObjectGuid string `xml:"OBJECTGUID,attr"`
	TypeName   string `xml:"TYPENAME,attr"`
	Level      int64  `xml:"LEVEL,attr"`
	Name       string `xml:"NAME,attr"`
	IsActive   bool   `xml:"ISACTIVE,attr"`
}

type House struct {
	ObjectId   int64  `xml:"OBJECTID,attr"`
	ObjectGuid string `xml:"OBJECTGUID,attr"`
	HouseNum   string `xml:"HOUSENUM,attr"`
	IsActive   bool   `xml:"ISACTIVE,attr"`
}

type Hierarchy struct {
	Id          int64 `xml:"ID,attr"`
	ObjectId    int64 `xml:"OBJECTID,attr"`
	ParentObjId int64 `xml:"PARENTOBJID,attr"`
	IsActive    bool  `xml:"ISACTIVE,attr"`
}

type Param struct {
	ObjectId  int64  `xml:"OBJECTID,attr"`
	TypeId    int64  `xml:"TYPEID,attr"`
	Value     string `xml:"VALUE,attr"`
	StartDate string `xml:"STARTDATE,attr"`
	EndDate   string `xml:"ENDDATE,attr"`
}

type AddressObjectType struct {
	Id        int64  `xml:"ID,attr"`
	Level     int64  `xml:"LEVEL,attr"`
	Name      string `xml:"NAME,attr"`
	ShortName string `xml:"SHORTNAME,attr"`
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
	List []FiasObject
}

func (a *FiasObjectList) AddObject(object FiasObject) {
	a.List = append(a.List, object)
}

func (a *FiasObjectList) Clear() {
	a.List = nil
}
