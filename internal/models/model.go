package models

type Model interface {
	Save() error
}

type ModelList interface {
	SaveModelList() error
}

type ModelStruct struct {
	tableName string
}
