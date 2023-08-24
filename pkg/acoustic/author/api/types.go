package api

func (element TextElement) Type() string {
	return "TextElement"
}

func (element MultiTextElement) Type() string {
	return "MultiTextElement"
}

func (element FormattedTextElement) Type() string {
	return "FormattedTextElement"
}

func (element NumberElement) Type() string {
	return "NumberElement"
}

func (element MultiNumberElement) Type() string {
	return "MultiNumberElement"
}

func (element FloatElement) Type() string {
	return "FloatElement"
}

func (element BooleanElement) Type() string {
	return "BooleanElement"
}

func (element LinkElement) Type() string {
	return "LinkElement"
}

func (element DateElement) Type() string {
	return "DateElement"
}

func (element CategoryElement) Type() string {
	return "CategoryElement"
}

func (element CategoryPartElement) Type() string {
	return "CategoryPartElement"
}

func (element ImageElement) Type() string {
	return "ImageElement"
}

func (element MultiImageElement) Type() string {
	return "MultiImageElement"
}

func (element FileElement) Type() string {
	return "FileElement"
}

func (element GroupElement) Type() string {
	return "GroupElement"
}

func (element MultiGroupElement) Type() string {
	return "MultiGroupElement"
}

func (element ReferenceElement) Type() string {
	return "ReferenceElement"
}

func (element MultiReferenceElement) Type() string {
	return "MultiReferenceElement"
}

func (element OptionSelectionElement) Type() string {
	return "OptionSelectionElement"
}

func (element MultiOptionSelectionElement) Type() string {
	return "MultiOptionSelectionElement"
}

func (element DateTimeElement) Type() string {
	return "DateTimeElement"
}

func (element MultiLinkElement) Type() string {
	return "MultiLinkElement"
}
