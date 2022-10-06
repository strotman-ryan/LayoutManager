package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func main() {
	router := gin.Default()
	router.GET("/component", getCustomComponentsGin)
	router.POST("/component", addCustomComponentGin)
	router.POST("/file", addFileGin)

	router.Run("localhost:8080")
}

// curl http://localhost:8080/component
func getCustomComponentsGin(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getCustomComponents())
}

//expecting
/*
{
	"componenName": {
		"key": "value1",
		"key1": "value",
		"key2": "value"
	}
}

example
{
	"Person": {
		"name": "STRING",
		"age": "FLOAT",
		"male": "BOOL"
	}
}

JSON good example
curl -X POST http://localhost:8080/component -H 'Content-Type: application/json' -d '{"Person": { "name": "STRING",  "age": "FLOAT", "male": "BOOL"}}'

JSON bad example
curl -X POST http://localhost:8080/component -H 'Content-Type: application/json' -d '{"Person": { "name": "STRING", , "age": "FLOAT"}}'
*/
func addCustomComponentGin(c *gin.Context) {
	var componentMap map[string]map[string]string
	if err := c.BindJSON(&componentMap); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//should only be one key, if there is more just return early
	for name, properties := range componentMap {
		if !addCustomComponent(name, properties) {
			c.String(http.StatusNotAcceptable, "Improper component properties")
			return
		}
		//only add the first one
		c.String(http.StatusOK, "All Good")
		return
	}
	c.String(http.StatusNotAcceptable, "Improper JSON format")
}

/*
Registers a file to an endpoint
expecting a JSON object
example 1:

	{
		"path": "/homepage",
		"type": "Person"
		"value": {
			"name": "Ryan",
			"age": 24,
			"male": true
		}
	}

curl -X POST http://localhost:8080/file -H 'Content-Type: application/json' -d '{"path": "/homepage", "type": "Person", "value": {"name": "Ryan", "age": 24,"male": true}}'
example 2:

	{
		"path": "/config",
		"type": "FLOAT"
		"value": 56
	}

curl -X POST http://localhost:8080/file -H 'Content-Type: application/json' -d '{"path": "/config", "type": "FLOAT", "value": 56}'
*/
func addFileGin(c *gin.Context) {
	var componentMap FileToJson
	if err := c.BindJSON(&componentMap); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if addFile(componentMap) {
		c.String(http.StatusOK, "All Good")
		return
	}
	c.String(http.StatusNotAcceptable, "Improper JSON format")
}

//---------------------NO GIN references below this line-------------------------

// returns false if add fails
func addFile(newJson FileToJson) bool {
	//need all these properties
	if newJson.Path == nil || newJson.ComponentType == nil || newJson.Value == nil {
		return false
	}

	//check the type matches the arbitrary data
	valid, errorString := jsonIsValid(*newJson.ComponentType, *newJson.Value)
	if !valid {
		fmt.Println(errorString)
	}

	//store the JSON
	if valid {
		//Will overwrite previous data
		endpoints[*newJson.Path] = *newJson.Value
	}
	fmt.Println(endpoints)
	return valid
}

// returns error stirng if json is not valid
func jsonIsValid(componentType string, value interface{}) (bool, string) {
	//check if it is a primitive component
	if function, exist := primitiveComponents[componentType]; exist {
		return function(value)
	}

	//check for customComponents
	componentProperyDefinitions, componentExists := customComponents[componentType]
	valueMap, valueExists := value.(map[string]interface{})
	if componentExists && valueExists {
		for propertyName, propertyType := range componentProperyDefinitions {
			value, exist := valueMap[propertyName]
			if !exist {
				//each property needs to be present
				//TODO: add optional support
				return false, fmt.Sprintf("Property: %v does not exist", propertyName)
			}
			if valid, error := jsonIsValid(propertyType, value); !valid {
				//JSON was not valid for that property
				return valid, error
			}
		}
		return true, "" //All properties were found and were valid
	}

	return false, fmt.Sprintf("Component Type: %v does not exist", componentType)
}

// I want to keep logic seperate from Gin incase of we switch frameworks
func getCustomComponents() map[string]map[string]string {
	return customComponents
}

// Will overwrite the component if one with that name already existed
// check if successful or not
func addCustomComponent(name string, properties map[string]string) bool {
	//check that all properties have existing types
	componentNames := getComponentNames()
	for _, value := range properties {
		//Must be a componentName already registered and no recursion is allowed
		if !slices.Contains(componentNames, value) || name == value {
			return false
		}
	}

	customComponents[name] = properties
	return true
}

//---------------------------- helper functions ----------------------------------

type FileToJson struct {
	//all types are references so can check null to see if they exist
	Path          *string      `json:"path"`
	ComponentType *string      `json:"type"`
	Value         *interface{} `json:"value"`
}

func getComponentNames() []string {
	//combine customComponents and builtInComponents names
	names := []string{}
	for name := range primitiveComponents {
		names = append(names, name)
	}

	for name := range customComponents {
		names = append(names, name)
	}

	return names
}

// returns error string if is not correct type
func isType[T interface{}](value interface{}) (bool, string) {
	_, ok := value.(T)
	if !ok {
		return ok, fmt.Sprintf("%v is of type %T; expecting type: %T", value, value, *new(T))
	}
	return ok, ""
}

// the key is the name of the component, value is a dictionary of its properties
//
//	the key is the name: the value is the type
var customComponents = make(map[string]map[string]string)

// like primitives
// this would be constant if go language allowed it
var primitiveComponents = map[string]func(interface{}) (bool, string){
	"FLOAT":  isType[float64],
	"STRING": isType[string],
	"BOOL":   isType[bool],
}

// The user's JSON
// key: the endpoint path like "/config" or "homepage/version3"
// value: raw JSON
var endpoints = map[string]interface{}{} //starts empty
