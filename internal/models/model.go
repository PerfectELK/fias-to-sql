package models

type Model interface {
	Save() error
}

type ModelList interface {
	SaveModelList() error
	AppendModel(mod Model)
}

type ModelStruct struct {
	tableName string
}
