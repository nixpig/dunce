package api

import (
	"fmt"
	"strconv"

	"github.com/nixpig/bloggor/internal/pkg/models"
)

func GetTags() map[string]models.TagData {
	tags, err := models.Query.Tag.GetAll()
	if err != nil {
		fmt.Println("error getting tags")
		return nil
	}

	tagmap := make(map[string]models.TagData)

	for index, item := range *tags {
		tagmap[strconv.Itoa(index)] = item
	}

	return tagmap
}
