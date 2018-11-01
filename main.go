package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Definition ...
type Definition struct {
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// Request ...
type Request struct {
	RqBody BodyObj `json:"rqBody"`
}

// Response ...
type Response struct {
	RsBody BodyObj   `json:"rsBody"`
	Error  BodyArray `json:"error"`
}

// FilterErrorCode ...
type FilterErrorCode struct {
	Code    BodyBasicType `json:"errorCode"`
	Message BodyBasicType `json:"errorDesc"`
}

// BodyObj ...
type BodyObj struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// BodyArray ...
type BodyArray struct {
	Type  string      `json:"type"`
	Items interface{} `json:"items"`
}

// BodyBasicType ...
type BodyBasicType struct {
	Type string `json:"type"`
}

var (
	dir = flag.String("dir", "swagger", "a directory")
)

func init() {
	flag.Parse()
}

func main() {

	filePath := *dir
	if filePath != "" {
		filePath += "/"
	}
	newDefinitions := make(map[string]interface{})
	JSONString, _ := ioutil.ReadFile(filePath + "swagger.json")
	var data map[string]interface{}
	json.Unmarshal(JSONString, &data)
	definitions := data["definitions"].(map[string]interface{})
	for i := range definitions {
		currentStruct := definitions[i].(map[string]interface{})
		tempProp := currentStruct["properties"]
		structName := strings.Split(i, ".")[1]
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(structName)), "req") {
			currentStruct["properties"] = Request{RqBody: BodyObj{Type: "object", Properties: tempProp}}
		} else if strings.HasPrefix(strings.ToLower(strings.TrimSpace(structName)), "res") {
			currentStruct["properties"] = Response{
				RsBody: BodyObj{Type: "object", Properties: tempProp},
				Error:  BodyArray{Type: "array", Items: BodyObj{Type: "object", Properties: FilterErrorCode{Code: BodyBasicType{Type: "string"}, Message: BodyBasicType{Type: "string"}}}},
			}
		} else {
			currentStruct["properties"] = tempProp
		}
		var newDef Definition
		mapstructure.Decode(currentStruct, &newDef)
		newDefinitions[i] = newDef
	}
	data["definitions"] = newDefinitions
	myJSON, _ := json.Marshal(data)

	f, err := os.Create(filePath + "swagger.json")
	if err != nil {
		log.Fatal("error create file", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(myJSON))
	if err != nil {
		log.Fatal("error write to "+f.Name()+".js", err)
		return
	}
	w.Flush()
}
