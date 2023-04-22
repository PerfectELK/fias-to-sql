package db

import (
	"fias_to_sql/internal/services/fias/types"
)

func ImportToDb(list *types.FiasObjectList) error {
	list.Clear()
	return nil
}
