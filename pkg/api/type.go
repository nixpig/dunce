package api

import (
	"fmt"
	"strconv"

	"github.com/nixpig/dunce/internal/pkg/models"
)

func GetTypes() map[string]models.Type {
	types, err := models.Query.Type.GetAll()
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get types: %v", err))
		return nil
	}

	typemap := make(map[string]models.Type)

	for index, item := range *types {
		typemap[strconv.Itoa(index)] = item
	}

	return typemap
}
