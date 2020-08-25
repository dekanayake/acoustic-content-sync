package domain

import "github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"

type Content struct {
	Name string `csv:"ContentName"`
	Type string `csv:"ContentType"`
	TypeId string `csv:"ContentTypeID"`
}

type Element struct {
	Name string        `csv:"ContentName"`
	Type api.FieldType `json:"type"`
}

type StringElement struct {
	Value string `json:"value"`
	Element
}

type LinkElement struct {
	LinkText string `json:"linkText"`
	Link string `json:"link"`
	Element
}

type CategoryElement struct {
	categories []string `json:"categories"`
	Element
}
