package models

import (
	"fias_to_sql/pkg/db"
	"fias_to_sql/pkg/db/types"
	"reflect"
	"strings"
)

type Model interface {
	Save() error
	GetTableName() string
	GetFields() []types.Key
	GetFieldValues() []string
}

type ModelList interface {
	SaveModelList() error
	AppendModel(mod Model)
}

type ModelListStruct struct {
	ModelList
	List []Model
}

func (r *ModelListStruct) AppendModel(mod Model) {
	r.List = append(r.List, mod)
}

type ModelStruct struct {
	TableName string
	Fields    []types.Key
}

func GetModelFields(m Model) []types.Key {
	ir := reflect.TypeOf(m)
	numFields := reflect.ValueOf(m).Elem().NumField()
	var keys []types.Key
	for i := 0; i < numFields; i++ {
		tag := ir.Elem().Field(i).Tag.Get("sql")
		if tag == "" {
			continue
		}
		tagArr := strings.Split(tag, ",")
		if tagArr[0] == "id" {
			continue
		}
		keys = append(keys, types.Key{
			Name: tagArr[0],
			Type: tagArr[1],
		})
	}
	return keys
}

func (r *ModelListStruct) SaveModelList() error {
	list := r.List
	if len(list) == 0 {
		return nil
	}
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	keys := list[0].GetFields()
	values := make([][]string, 0, len(list))
	for _, item := range list {
		fieldValues := item.GetFieldValues()
		values = append(values, fieldValues)
	}

	tableName := reflect.Indirect(reflect.ValueOf(list[0])).FieldByName("TableName").String()
	err = dbInstance.InsertList(tableName, keys, values)
	return err
}
