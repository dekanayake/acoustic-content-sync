package api

type GenericData struct {
	Name string
	Type string
	Value string
}


func TextConverter (data interface{},element Element) (Element,error) {
	textElement := element.(TextElement)
	textElement.Value = data.(GenericData).Value
	return textElement,nil
}
