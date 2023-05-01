package models

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
	"strconv"
)

type Hierarchy struct {
	Model
	ModelStruct
	id               int64 `sql:"id,int"`
	object_id        int64 `sql:"object_id,int"`
	parent_object_id int64 `sql:"parent_object_id,int"`
}

func (h *Hierarchy) SetId(id int64) {
	h.id = id
}

func (h *Hierarchy) SetObject_id(object_id int64) {
	h.object_id = object_id
}

func (h *Hierarchy) SetParent_object_id(parent_object_id int64) {
	h.parent_object_id = parent_object_id
}

func NewHierarchy() *Hierarchy {
	object := Hierarchy{}
	object.TableName = config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	return &object
}

func (h *Hierarchy) Save() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	queryMap := map[string]string{
		"object_id":        strconv.FormatInt(h.object_id, 10),
		"parent_object_id": strconv.FormatInt(h.parent_object_id, 10),
	}
	return dbInstance.Insert(h.TableName, queryMap)
}
