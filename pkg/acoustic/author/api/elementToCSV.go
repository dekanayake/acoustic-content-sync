package api

import (
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	errors "github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"strconv"
	"strings"
)

func (element NumberElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{
		Value: strconv.Itoa(int(element.Value)),
	}, nil
}

func (element TextElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{
		Value: element.Value,
	}, nil
}

func (element LinkElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element FormattedTextElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element MultiTextElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element FloatElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{Value: fmt.Sprintf("%.2f", element.Value)}, nil
}

func (element BooleanElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element DateElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element CategoryElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	categories := element.Categories
	outputCats := make([]string, 0)
	for _, cat := range categories {
		splitCats := strings.Split(cat, "/")
		splitCats = splitCats[1:]
		outputCats = append(outputCats, strings.Join(splitCats, env.CategoryHierarchySeperator()))
	}
	return CSVValues{
		Value: strings.Join(outputCats, env.MultipleItemsSeperator()),
	}, nil
}

func (element CategoryPartElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element ImageElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	assetID := element.Asset.ID
	response, err := NewAssetClient(env.AcousticAPIUrl()).Get(assetID)
	if err != nil {
		errors.ErrorWithStack(err)
	}
	return CSVValues{
		Value: env.AcousticDomain() + response.Path,
	}, nil
}

func (element FileElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element GroupElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	csvValues := make(map[string]CSVValues, 0)
	values := element.Value
	for key, value := range values {
		childFieldMapping := fieldMapping.GetChildFieldMapping(key)
		if childFieldMapping != nil {
			childCsvValues, err := value.(Element).ToCSV(childFieldMapping)
			if err != nil {
				return CSVValues{}, errors.ErrorWithStack(err)
			}
			csvValues[key] = childCsvValues
		}
	}
	return CSVValues{
		ChildValues: csvValues,
	}, nil

}

func (element MultiGroupElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element ReferenceElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element MultiReferenceElement) ToCSV(fieldMapping *csv.ContentFieldMapping) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}
