package api

import (
	"github.com/wesovilabs/koazee"
	"strconv"
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


