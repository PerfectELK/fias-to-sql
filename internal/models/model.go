package models

import (
	"fias_to_sql/pkg/db"
	"fias_to_sql/pkg/db/types"
	"reflect"
	"strconv"
	"strings"
)

type Model interface {
	Save() error
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
}

func (r *ModelListStruct) SaveModelList() error {
	list := r.List
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	var keys []types.Key
	if len(list) == 0 {
		return nil
	}
	ir := reflect.TypeOf(list[0])
	numFields := reflect.ValueOf(list[0]).Elem().NumField()
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

	var values [][]string
	for _, item := range list {
		r := reflect.ValueOf(item)
		var value []string
		for _, key := range keys {
			v := reflect.Indirect(r).FieldByName(key.Name)
			if key.Type == "string" {
				value = append(value, v.String())
			}
			if key.Type == "int" {
				value = append(value, strconv.FormatInt(v.Int(), 10))
			}
		}
		values = append(values, value)
	}

	tableName := reflect.Indirect(reflect.ValueOf(list[0])).FieldByName("TableName").String()
	return dbInstance.InsertList(tableName, keys, values)
}
