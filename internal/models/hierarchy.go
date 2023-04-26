package models

import (
	"fias_to_sql/pkg/db"
	"strconv"
)

type Hierarchy struct {
	Model
	ModelStruct
	id               int64
	object_id        int64
	parent_object_id int64
}

type HierarchyList struct {
	ModelList
	List []Model
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
	object.tableName = "fias_objects_hierarchy"
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
	return dbInstance.Insert(h.tableName, queryMap)
}

func (m *HierarchyList) SaveModelList() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	keys := []string{
		"object_id",
		"parent_object_id",
	}
	tableName := ""
	var values [][]string
	for _, val := range m.List {
		objVal, _ := val.(*Hierarchy)
		if tableName == "" {
			tableName = objVal.tableName
		}
		vals := []string{
			strconv.FormatInt(objVal.object_id, 10),
			strconv.FormatInt(objVal.parent_object_id, 10),
		}
		values = append(values, vals)
	}

	return dbInstance.InsertList(tableName, keys, values)
}

func (m *HierarchyList) AppendModel(mod Model) {
	m.List = append(m.List, mod)
}
