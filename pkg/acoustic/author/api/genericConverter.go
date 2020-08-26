package api

import (
	"errors"
	"github.com/dekanayake/acoustic-content-sync/pkg/context"
	"github.com/wesovilabs/koazee"
	"strconv"
	"strings"
)

type AcousticDataRecord struct {
	NameFields []string
	Values []GenericData
}

type GenericData struct {
	Name string
	Type string
	Value string
}

func (acousticDataRecord AcousticDataRecord) Name() string {
	values := make([]GenericData,0,len(acousticDataRecord.NameFields))
	for _,nameField := range acousticDataRecord.NameFields {
		value:=koazee.StreamOf(acousticDataRecord.Values).
			Filter(func(data GenericData) bool{
				return data.Name == nameField
		}).First().Val().(GenericData)
		values = append(values, value)
	}

	return koazee.StreamOf(values).
		Reduce(func(acc string, data GenericData) string{
			if acc == "" {
				acc += data.Value
			} else {
				acc += "__" + data.Value
			}
			return acc
	}).String()
}


func (element TextElement) Convert (data interface{}) (Element,error) {
	element.Value = data.(GenericData).Value
	return element,nil
}

func (element NumberElement) Convert (data interface{}) (Element,error) {
	numValue,err := strconv.ParseInt(data.(GenericData).Value,0, 64)
	if err != nil {
		return nil,err
	}
	element.Value = numValue
	return element,nil
}


func (element LinkElement) Convert (data interface{}) (Element,error) {
	element.LinkURL = data.(GenericData).Value
	return element,nil
}


func (element CategoryElement) Convert (data interface{}) (Element,error) {
	catItems := strings.Split(data.(GenericData).Value,"/")
	if len(catItems) == 1 {
		return nil, errors.New("empty category :" + catItems[0])
	}

	categoryItems,err := NewCategoryClient(context.AcousticAPIUrl()).Categories(catItems[0])
	if err != nil {
		return nil, err
	}
	catNamePaths := make([]string,0,10)
	for i := 1; i <= len(catItems); i++ {
		catNamePathsSlice := catItems[0:i]
		catNamePath := koazee.StreamOf(catNamePathsSlice).
			Reduce(func(acc string, catNamePath string) string {
				if acc == "" {
					acc += catNamePath
				} else {
					acc += "/" + catNamePath
				}
				return acc
		}).String()
		catNamePaths = append(catNamePaths, catNamePath)
	}

	catIds := koazee.StreamOf(categoryItems).
		Filter(func(categoryItem CategoryItem) bool{
			fullNamePath := categoryItem.FullNamePath()
			contains,_ :=  koazee.StreamOf(catNamePaths).Contains(fullNamePath)
			return strings.Contains(fullNamePath,"/") && contains
	}).
		Map(func(categoryItem CategoryItem) string {
			return categoryItem.Id
	}).Out().Val().([]string)

	element.CategoryIds = catIds
	return element,nil
}


