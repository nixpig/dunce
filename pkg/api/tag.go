package api

import (
	"fmt"
	"strconv"

	"github.com/nixpig/dunce/internal/pkg/models"
)

func GetTags() map[string]models.Tag {
	tags, err := models.Query.Tag.GetAll()
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get tags: %v", err))
		return nil
	}

	tagmap := make(map[string]models.Tag)

	for index, item := range *tags {
		tagmap[strconv.Itoa(index)] = item
	}

	return tagmap
}
