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
		"key": "value1"
		"key1": "value"
		"key2": "value"
	}
}

example
{
	"Person": {
		"name": "STRING"
		"age": "INT"
		"male": "BOOL"
	}
}

JSON good example
curl -X POST http://localhost:8080/component -H 'Content-Type: application/json' -d '{"Person": { "name": "STRING",  "age": "INT"}}'

JSON bad example
curl -X POST http://localhost:8080/component -H 'Content-Type: application/json' -d '{"Person": { "name": "STRING", , "age": "INT"}}'
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
example:

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
example:

	{
		"path": "/config",
		"type": "INT"
		"value": 56
	}
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
		fmt.Println("expected value not present")
		return false
	}

	//check componetType is already registered
	if !slices.Contains(getComponentNames(), *(newJson.ComponentType)) {
		fmt.Println("componentType: not registered")
		return false
	}

	//check the type matches the arbitrary data
	return true
}

func jsonIsValid(componentType string, value interface{}) bool {
	return false //TODO: Implement
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
	for _, name := range primitiveComponents {
		names = append(names, name)
	}

	for name := range customComponents {
		names = append(names, name)
	}

	return names
}

// the key is the name of the component, value is a dictionary of its properties
//
//	the key is the name: the value is the type
var customComponents = make(map[string]map[string]string)

// like primitives
// this would be constant if go language allowed it
var primitiveComponents = [3]string{"INT", "STRING", "BOOL"}
