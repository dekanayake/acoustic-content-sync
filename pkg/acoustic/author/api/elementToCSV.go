package api

import (
	"errors"
	"strconv"
)

func (element NumberElement) ToCSV() (string, error) {
	return strconv.Itoa(int(element.Value)), nil
}

func (element TextElement) ToCSV() (string, error) {
	return element.Value, nil
}

func (element LinkElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element FormattedTextElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element MultiTextElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element FloatElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element BooleanElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element DateElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element CategoryElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element CategoryPartElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element ImageElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element FileElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element GroupElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element MultiGroupElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element ReferenceElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}

func (element MultiReferenceElement) ToCSV() (string, error) {
	return "", errors.New("to csv not implemented")
}
