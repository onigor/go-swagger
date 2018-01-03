package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	UndefinedStructType = "q!undefined"
)

type SwaggerMainDocStruct struct {
	Swagger      string                       `json:"swagger"`
	Info         SwaggerInfoStruct            `json:"info"`
	Definitions  map[string]SwaggerDocStruct  `json:"definitions"`
	Host         string                       `json:"host"`
	BasePath     string                       `json:"basePath"`
	Tags         []SwaggerTagStruct           `json:"tags"`
	Schemes      []string                     `json:"schemes"`
	Paths        map[string]SwaggerPathStruct `json:"paths"`
	ExternalDocs SwaggerExternalDocs          `json:"externalDocs"`
}

type SwaggerPathStruct struct {
}

type SwaggerTagStruct struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ExternalDocs SwaggerExternalDocs `json:"externalDocs"`
}

type SwaggerExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type SwaggerInfoStruct struct {
	Description string               `json:"description"`
	Version     string               `json:"version"`
	Title       string               `json:"title"`
	Contact     SwaggerContactStruct `json:"contact"`
	License     SwaggerLicenseStruct `json:"license"`
}

type SwaggerLicenseStruct struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SwaggerContactStruct struct {
	Email string `json:"email"`
}

type SwaggerDocStruct struct {
	Extend     *[]string                     `json:"extend,omitempty"`
	Type       string                        `json:"type"`
	Name       string                        `json:"name"`
	Properties map[string]SwaggerDocProperty `json:"properties"`
	XML        SwaggerDocXMLStruct           `json:"xml,omitempty"`
}

type SwaggerDocXMLStruct struct {
	Name string `json:"name,omitempty"`
}

func (s SwaggerDocStruct) MarshalJSON() ([]byte, error) {
	log("in:SwaggerDocStruct")
	type Alias SwaggerDocStruct
	return json.Marshal(&struct {
		Type string `json:"type"`
		Name string `json:"name,omitempty"`
		*Alias
	}{
		Type:  correctType(s.Type),
		Name:  "",
		Alias: (*Alias)(&s),
	})
}

func correctType(str string) string {
	if str == "interface" {
		return "object"
	}
	str = strings.Replace(str, "struct", "object", -1)
	//todo map[string]{something}
	str = strings.Replace(str, "map[string]interface{}", "object", -1)
	str = strings.Replace(str, "interface{}", "object", -1)
	if strings.Contains(str, "int32") || strings.Contains(str, "int64") || str == "int" {
		return "integer"
	}
	return str
}

type SwaggerDocProperty struct {
	Type        string        `json:"type,omitempty"`
	Ref         string        `json:"$ref,omitempty"`
	Description string        `json:"description,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	Format      string        `json:"format,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
}

func (s SwaggerDocProperty) MarshalJSON() ([]byte, error) {
	log("in:SwaggerDocProperty", s.Type)
	log("in:SwaggerDocProperty $ref", s.Ref)

	if len(s.Ref) != 0 {
		return json.Marshal(&struct {
			Ref string `json:"$ref,omitempty"`
		}{
			Ref: s.Ref,
		})
	}

	type Alias SwaggerDocProperty
	return json.Marshal(&struct {
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
		*Alias
	}{
		Type:  correctType(s.Type),
		Name:  "",
		Alias: (*Alias)(&s),
	})
}

const (
	OutputFileMessage = "Swagger JSON Spec output file"
	DebugMessage      = "Debug mode, print some steps"
)

var OutputFile string
var DebugMode bool

func init() {
	flag.StringVar(&OutputFile, "o", "swagger.json", OutputFileMessage)
	flag.BoolVar(&DebugMode, "d", false, DebugMessage)
	// flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")
}

func log(args ...interface{}) {
	if DebugMode {
		fmt.Println(args...)
	}
}

const (
	TestFilesExcludeSuffix = "_test.go"
)

func main() {
	flag.Parse()

	ext := `.go`

	log("Hello")
	log("Program is in debug mode")

	swg := SwaggerMainDocStruct{}

	basePath, _ := os.Getwd()
	configData, err := ioutil.ReadFile(filepath.Join(basePath, ".swagit.json"))
	if err == nil {
		json.Unmarshal(configData, &swg)
	}

	filePaths := checkExt(ext)
	totalList := map[string]SwaggerDocStruct{}

	if len(filePaths) > 0 {
		for index := 0; index < len(filePaths); index++ {
			list, err := parseFile(filePaths[index])
			if err != nil {
				panic(err)
			}
			if len(list) > 0 {
				for index := 0; index < len(list); index++ {
					totalList[list[index].Name] = list[index]
				}
				// totalList = append(totalList, list...)
			}
		}
	} else {
		log("No files found with extension ", ext)
		return
	}
	swg.Definitions = totalList

	data, err := json.MarshalIndent(swg, "", "  ")
	if err != nil {
		panic(err)
	}

	finalOutputFilePath := filepath.Join(basePath, OutputFile)
	err = ioutil.WriteFile(finalOutputFilePath, data, os.ModeAppend)
	if err != nil {
		log("Error writing to file", finalOutputFilePath)
		fmt.Println(string(data))
	} else {
		fmt.Println("Spec file was successfully generated", finalOutputFilePath)
	}
	log("The end")
}

func parseFile(fPath string) ([]SwaggerDocStruct, error) {
	result := []SwaggerDocStruct{}
	log("File path ", fPath)

	packagePrefix := ""

	// validPackageName := regexp.MustCompile(`package\s+(\w+)`)

	validGoStruct := regexp.MustCompile(`type\s+(\w+)\s+(\w+)\s*(?:{((?:[^{]|(?:{}))*)})`)
	fileData, err := ioutil.ReadFile(fPath)
	if err != nil {
		return result, err
	}
	fileDataString := string(fileData)
	// packagePrefixResults := validPackageName.FindStringSubmatch(fileDataString)
	// if len(packagePrefixResults) == 2 {
	// 	packagePrefix = packagePrefixResults[1]
	// }
	foundArr := validGoStruct.FindAllStringSubmatch(fileDataString, -1)
	if len(foundArr) == 0 {
		return result, nil
	}

	for index := 0; index < len(foundArr); index++ {
		obj, err := parseStruct(foundArr[index], packagePrefix)
		if err != nil {
			return result, err
		}
		result = append(result, obj)
	}

	return result, nil
}

func parseStruct(data []string, packagePrefix string) (SwaggerDocStruct, error) {
	// fmt.Printf("len=%d,data=%+s", len(data), )
	for _, item := range data {
		log("item part=", strings.TrimSpace(string(item)))
	}
	if len(data) < 4 {
		return SwaggerDocStruct{}, fmt.Errorf("parse error data=%s", data)
	}

	// structFull := data[0]
	structName := data[1]
	structType := data[2]
	structBody := data[3]

	result := SwaggerDocStruct{}
	result.Type = string(structType)
	result.Name = packagePrefix + string(structName)
	result.XML = SwaggerDocXMLStruct{Name: result.Name}

	if len(structBody) == 0 {
		log("Error!", "empty struct")
		return result, nil //fmt.Errorf("empty struct")
	}

	structBodyString := trimString(fixNewLine(structBody))
	lines := strings.Split(structBodyString, "\n")

	if len(lines) == 0 {
		log("Error!", "empty struct 2")
		return result, nil //fmt.Errorf("empty struct 2")
	}

	if len(packagePrefix) != 0 {
		packagePrefix = packagePrefix + "."
	}

	result.Properties = map[string]SwaggerDocProperty{}

	for index := 0; index < len(lines); index++ {
		line := trimString(lines[index])
		if len(line) == 0 {
			continue
		}
		lineParts := strings.Split(line, " ")
		if len(lineParts) == 1 {
			// result.Extend = append(result.Extend, lineParts[0])
			// class extention
			//load class now or later on sereliazation
		}

		if len(lineParts) == 2 {
			swdp := SwaggerDocProperty{}
			swdp.Type = lineParts[1]
			if customType(swdp.Type) {
				swdp.Ref = "#/definitions/" + swdp.Type
			}
			result.Properties[lineParts[0]] = swdp
			//normal property
		}

		if len(lineParts) > 2 {
			key := parsePropertyKeyName(line, lineParts[0])
			swdp := SwaggerDocProperty{}
			swdp.Type = lineParts[1]
			if customType(swdp.Type) {
				swdp.Ref = "#/definitions/" + swdp.Type
			}
			swdp.Description += extractComments(line)

			result.Properties[key] = swdp
			//extended property
		}
	}

	log("____")
	return result, nil
}

func inArray(arr []string, value string) bool {
	for index := 0; index < len(arr); index++ {
		if strings.Contains(value, arr[index]) {
			return true
		}
	}
	return false
}

func customType(str string) bool {
	excludeArr := []string{"int", "int32", "int64", "float", "string", "struct", "interface{}", "bool", "map[string]interface{}", "[]interface{}"}
	return !inArray(excludeArr, str)
}

func parsePropertyKeyName(str string, def string) string {
	validGoNotation := regexp.MustCompile("`.*`")
	jsonPropertyName := regexp.MustCompile(`json\:\"(.*)\"`)
	if foundStr := validGoNotation.FindString(str); len(foundStr) > 0 {
		results := jsonPropertyName.FindStringSubmatch(str)
		if len(results) > 1 {
			splitted := strings.Split(results[1], ",")
			return splitted[0]
		}
	}
	return def
}

func extractComments(str string) string {
	if strings.Contains(str, "//") {

		return str[strings.Index(str, "//"):] + " "
	}
	return ""
}

func fixNewLine(str string) string {
	str = strings.Replace(str, "\n\r", "\n", -1)
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "\n", -1)
	return str
}

func trimString(str string) string {
	emptySpace := regexp.MustCompile(`\t+| {2,}`)
	// str = strings.Replace(str, "\t", " ", -1)
	str = emptySpace.ReplaceAllString(str, " ")
	str = strings.Replace(str, "  ", " ", -1)
	str = strings.TrimPrefix(str, "\n")
	str = strings.TrimSuffix(str, "\n")
	return strings.TrimSpace(str)
}

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		log("Walk path", path)
		if !f.IsDir() {
			if strings.HasSuffix(path, ext) && !strings.HasSuffix(path, TestFilesExcludeSuffix) {
				log("selected", path)
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}
