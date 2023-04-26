package models

import (
	"fias_to_sql/pkg/db"
	"strconv"
)

type Object struct {
	Model
	ModelStruct
	id          int64
	object_id   int64
	object_guid string
	type_name   string
	level       int64
	name        string
	add_name    string
	add_name2   string
}

type ObjectList struct {
	ModelList
	List []Model
}

func (m *Object) SetObject_id(object_id int64) {
	m.object_id = object_id
}

func (m *Object) SetObject_guid(object_guid string) {
	m.object_guid = object_guid
}

func (m *Object) SetType_name(type_name string) {
	m.type_name = type_name
}

func (m *Object) SetLevel(level int64) {
	m.level = level
}

func (m *Object) SetName(name string) {
	m.name = name
}

func (m *Object) SetAdd_name(add_name string) {
	m.add_name = add_name
}

func (m *Object) SetAdd_name2(add_name2 string) {
	m.add_name2 = add_name2
}

func NewObject() *Object {
	object := Object{}
	object.tableName = "fias_objects"
	return &object
}

func (m *Object) Save() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	queryMap := map[string]string{
		"object_id":   strconv.FormatInt(m.object_id, 10),
		"object_guid": m.object_guid,
		"type_name":   m.type_name,
		"level":       strconv.FormatInt(m.level, 10),
		"name":        m.name,
		"add_name":    m.add_name,
		"add_name2":   m.add_name2,
	}
	return dbInstance.Insert(m.tableName, queryMap)
}

func (m *ObjectList) SaveModelList() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	keys := []string{
		"object_id",
		"object_guid",
		"type_name",
		"level",
		"name",
		"add_name",
		"add_name2",
	}
	tableName := ""
	var values [][]string
	for _, val := range m.List {
		objVal, _ := val.(*Object)
		if tableName == "" {
			tableName = objVal.tableName
		}
		vals := []string{
			strconv.FormatInt(objVal.object_id, 10),
			objVal.object_guid,
			objVal.type_name,
			strconv.FormatInt(objVal.level, 10),
			objVal.name,
			objVal.add_name,
			objVal.add_name2,
		}
		values = append(values, vals)
	}

	return dbInstance.InsertList(tableName, keys, values)
}

func (m *ObjectList) AppendModel(mod Model) {
	m.List = append(m.List, mod)
}
