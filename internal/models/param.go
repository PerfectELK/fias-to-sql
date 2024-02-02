package models

import (
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/pkg/db"
	"github.com/PerfectELK/go-import-fias/pkg/db/types"
	"strconv"
)

type Param struct {
	ModelStruct
	id        int64  `sql:"id,int"`
	object_id int64  `sql:"object_id,int"`
	kladr_id  string `sql:"kladr_id,string"`
}

var paramFields []types.Key

func (h *Param) SetId(id int64) {
	h.id = id
}

func (h *Param) SetObject_id(object_id int64) {
	h.object_id = object_id
}

func (h *Param) SetKladr_id(kladr_id string) {
	h.kladr_id = kladr_id
}

func (h *Param) GetTableName() string {
	return h.TableName
}

func NewParam() *Param {
	object := Param{}
	object.TableName = config.GetConfig("DB_OBJECTS_KLADR_TABLE")
	return &object
}

func (h *Param) Save() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	queryMap := map[string]string{
		"object_id": strconv.FormatInt(h.object_id, 10),
		"kladr_id":  h.kladr_id,
	}
	return dbInstance.Insert(h.TableName, queryMap)
}

func (m *Param) GetFields() []types.Key {
	if len(paramFields) != 0 {
		return paramFields
	}
	paramFields = GetModelFields(m)
	return paramFields
}

func (m *Param) GetFieldValues() []string {
	return []string{
		strconv.FormatInt(m.object_id, 10),
		m.kladr_id,
	}
}
