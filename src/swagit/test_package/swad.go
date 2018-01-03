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

type lowerCaseTestStruct struct {
	loverCasePropertyName string
}

type Duper struct {
	SuperDuperEmptyStruct
	SuperDuperStruct
	DuperProp string `json:"duperPropOmit,omitempty"`
}

type NotEmpty struct {
	SuperDuperStruct
}

type TestArrayStruct struct {
	Test                    []string
	TestingMote             []Duper
	TestEvenMoreWithPointer []*Duper
	TestPointer             *[]Duper
	TestStringPointer       *[]string
}

type All struct {
	Super
	SuperDuperEmptyStruct
	loverCasePropertyName string
	SuperProp             string
	DuperProp             string `json:"duperProp"`
	NotEmptyProp          int    //this is int property
}
