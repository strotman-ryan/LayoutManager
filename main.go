package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func main() {
	router := gin.Default()
	router.GET("/component", getCustomComponentsGin)
	router.POST("/component", addCustomComponentGin)

	router.Run("localhost:8080")
}

// curl http://localhost:8080/component
func getCustomComponentsGin(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getCustomComponents())
}

//expecting
/*
component = Serialized JSON
"componenName": {
	"parameter1name": "value1"
	"key1": value
	"key2": value
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
	return
}

//---------------------NO GIN references below this line-------------------------

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

func getComponentNames() []string {
	//combine customComponents and builtInComponents names
	names := []string{}
	for _, name := range builtInComponents {
		names = append(names, name)
	}

	for name, _ := range customComponents {
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
var builtInComponents = [3]string{"INT", "STRING", "BOOL"}
