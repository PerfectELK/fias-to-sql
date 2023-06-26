package types

const (
	IMPORT_DUMP_FILENAME = "dump.json"
)

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
	IsActive   bool   `xml:"ISACTIVE,attr"`
}

type Hierarchy struct {
	FiasObject
	Id          int64 `xml:"ID,attr"`
	ObjectId    int64 `xml:"OBJECTID,attr"`
	ParentObjId int64 `xml:"PARENTOBJID,attr"`
	IsActive    bool  `xml:"ISACTIVE,attr"`
}

type Param struct {
	FiasObject
	ObjectId  int64  `xml:"OBJECTID,attr"`
	TypeId    int64  `xml:"TYPEID,attr"`
	Value     string `xml:"VALUE,attr"`
	StartDate string `xml:"STARTDATE,attr"`
	EndDate   string `xml:"ENDDATE,attr"`
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

type Dump struct {
	ArchivePath string   `json:"archive_path"`
	TablesType  string   `json:"tables_type"`
	Files       []string `json:"files"`
}

func MakeDump() error {
	return nil
}
