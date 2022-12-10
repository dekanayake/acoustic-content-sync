package api

import (
	"fmt"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	errors "github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"strconv"
	"strings"
)

func (element NumberElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{
		Value: strconv.Itoa(int(element.Value)),
	}, nil
}

func (element TextElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{
		Value: element.Value,
	}, nil
}

func (element LinkElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{
		Value: element.LinkURL,
	}, nil
}

func (element FormattedTextElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element MultiTextElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element FloatElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{Value: fmt.Sprintf("%.6f", element.Value)}, nil
}

func (element BooleanElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{
		Value: strconv.FormatBool(element.Value),
	}, nil
}

func (element DateElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element CategoryElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
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

func (element CategoryPartElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element ImageElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	url := ""
	if element.URL != "" {
		url = env.AcousticDomain() + element.URL
	} else {
		assetID := element.Asset.ID
		response, err := NewAssetClient(env.AcousticAPIUrl()).Get(assetID)
		if err != nil {
			errors.ErrorWithStack(err)
		}
		url = env.AcousticDomain() + response.Path
	}

	return CSVValues{
		Value: url,
	}, nil
}

func (element FileElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element GroupElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	csvValues := make(map[string]CSVValues, 0)
	values := element.Value
	for key, value := range values {
		nextLevelChildFields := childFields[key]
		if nextLevelChildFields != nil {
			childCsvValues, err := value.(Element).ToCSV(nextLevelChildFields.(map[string]interface{}))
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

func (element MultiGroupElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element ReferenceElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (element MultiReferenceElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	return CSVValues{}, errors.ErrorMessageWithStack("to csv not implemented")
}

func (o OptionSelectionElement) ToCSV(childFields map[string]interface{}) (CSVValues, error) {
	//TODO implement me
	panic("implement me")
}
