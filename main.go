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

// MyRequest ...
type MyRequest struct {
	RqBody MyBodyObj `json:"rqBody"`
}

// MyResponse ...
type MyResponse struct {
	RsBody MyBodyObj `json:"rsBody"`
}

// MyBodyObj ...
type MyBodyObj struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

var (
	dir = flag.String("dir", "", "a directory")
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
			currentStruct["properties"] = MyRequest{RqBody: MyBodyObj{Type: "object", Properties: tempProp}}
		} else if strings.HasPrefix(strings.ToLower(strings.TrimSpace(structName)), "res") {
			currentStruct["properties"] = MyResponse{RsBody: MyBodyObj{Type: "object", Properties: tempProp}}
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
