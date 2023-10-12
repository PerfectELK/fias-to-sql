package models

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
	"fias_to_sql/pkg/db/types"
	"strconv"
)

type ObjectType struct {
	Model
	ModelStruct
	id         int64  `sql:"id,int"`
	level      int64  `sql:"level,int"`
	name       string `sql:"name,string"`
	short_name string `sql:"short_name,string"`
}

func (o *ObjectType) SetId(id int64) {
	o.id = id
}

func (o *ObjectType) SetLevel(level int64) {
	o.level = level
}

func (o *ObjectType) SetName(name string) {
	o.name = name
}

func (o *ObjectType) SetShortName(short_name string) {
	o.short_name = short_name
}

func NewObjectType() *ObjectType {
	object := ObjectType{}
	object.TableName = config.GetConfig("DB_OBJECT_TYPES_TABLE")
	return &object
}

func (m *ObjectType) Save() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	queryMap := map[string]string{
		"id":         strconv.FormatInt(m.id, 10),
		"level":      strconv.FormatInt(m.level, 10),
		"name":       m.name,
		"short_name": m.short_name,
	}
	return dbInstance.Insert(m.TableName, queryMap)
}

func (m *ObjectType) GetFields() []types.Key {
	if len(m.Fields) != 0 {
		return m.Fields
	}
	m.Fields = GetModelFields(m)
	return m.Fields
}

func (m *ObjectType) GetFieldValues() []string {
	return []string{
		strconv.FormatInt(m.level, 10),
		m.name,
		m.short_name,
	}
}