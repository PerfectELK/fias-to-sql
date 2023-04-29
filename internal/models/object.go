package models

import (
	"fias_to_sql/pkg/db"
	"strconv"
)

type Object struct {
	Model
	ModelStruct
	id          int64  `sql:"id,int"`
	object_id   int64  `sql:"object_id,int"`
	object_guid string `sql:"object_guid,string"`
	type_name   string `sql:"type_name,string"`
	level       int64  `sql:"level,int"`
	name        string `sql:"name,string"`
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

func NewObject() *Object {
	object := Object{}
	object.TableName = "fias_objects"
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
	}
	return dbInstance.Insert(m.TableName, queryMap)
}
