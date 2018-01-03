package custPackage

type SuperDuperEmptyStruct struct {
}

type SuperDuperStruct struct {
	Super    string
	Duper    string `json:"duper"`
	NotEmpty Duper  `json:"customName"` //this one is just for fun
}

type Super struct {
	SuperDuperStruct
}

type SuperInterface interface {
}

type Duper struct {
	SuperDuperEmptyStruct
	SuperDuperStruct
	DuperProp string `json:"duperPropOmit,omitempty"`
}

type NotEmpty struct {
	SuperDuperStruct
}

type All struct {
	Super
	SuperDuperEmptyStruct
	SuperProp    string
	DuperProp    string `json:"duperProp"`
	NotEmptyProp int    //this is int property
}
