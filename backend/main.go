package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

const arrayName = "ARRAY"

func main() {
	//initlizing `primitiveComponents` here to break an initlization cycle
	primitiveComponents = map[string]func(interface{}) (bool, interface{}){
		"FLOAT":   isType[float64],
		"INT":     isInt,
		"STRING":  isType[string],
		"BOOL":    isType[bool],
		arrayName: isArray,
	}
	router := gin.Default()

	//Set up Cors
	//TODO SECURITY RISK: make this more granular
	//This allows all orgins
	router.Use(cors.Default())

	router.GET("/component", getCustomComponentsGin)
	router.POST("/component", addCustomComponentGin)
	router.POST("/file", addFileGin)
	router.GET("/:urlLocation", getJSON)

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
			c.String(http.StatusBadRequest, "Improper component properties")
			return
		}
		//only add the first one
		c.String(http.StatusOK, "All Good")
		return
	}
	c.String(http.StatusBadRequest, "Improper JSON format")
}

/*
Registers a file to an endpoint
expecting a JSON object
path: should NOT start with a "/"
example 1:

	{
		"path": "homepage",
		"type": "Person"
		"value": {
			"name": "Ryan",
			"age": 24,
			"male": true
		}
	}

curl -X POST http://localhost:8080/file -H 'Content-Type: application/json' -d '{"path": "homepage", "type": "Person", "value": {"name": "Ryan", "age": 24,"male": true}}'
example 2:

	{
		"path": "config",
		"type": "FLOAT"
		"value": 56
	}

curl -X POST http://localhost:8080/file -H 'Content-Type: application/json' -d '{"path": "config", "type": "FLOAT", "value": 56}'
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
	c.String(http.StatusBadRequest, "Improper JSON format")
}

/* Get the raw JSON the user placed on the server
curl http://localhost:8080/homepage
curl http://localhost:8080/config
*/

func getJSON(c *gin.Context) {
	urlPath := c.Param("urlLocation")
	if file, exist := files[urlPath]; exist {
		//get the raw value; The value should always be valid at this point
		if valid, rawValue := jsonIsValid(*file.ComponentType, *file.Value); valid {
			c.JSON(http.StatusOK, rawValue)
			return
		}
	}
	c.String(http.StatusBadRequest, "JSON does not exist")
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
		files[*newJson.Path] = newJson
	}
	return valid
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

type ArrayItem struct {
	ComponentType *string      `mapstructure:"type"`
	Value         *interface{} `mapstructure:"value"`
}

// returns error stirng if json is not valid
// if valid it returns the raw JSON form
func jsonIsValid(componentType string, value interface{}) (bool, interface{}) {
	//check if it is a primitive component
	if function, exist := primitiveComponents[componentType]; exist {
		return function(value)
	}

	//check for customComponents
	componentProperyDefinitions, componentExists := customComponents[componentType]
	valueMap, valueExists := value.(map[string]interface{})
	if componentExists && valueExists {
		var rawValues = make(map[string]interface{})
		for propertyName, propertyType := range componentProperyDefinitions {
			value, exist := valueMap[propertyName]
			if !exist {
				//each property needs to be present
				//TODO: add optional support
				return false, fmt.Sprintf("Property: %v does not exist", propertyName)
			}
			valid, rawValue := jsonIsValid(propertyType, value)
			if valid {
				rawValues[propertyName] = rawValue
			} else {
				//JSON was not valid for that property
				return valid, rawValue //rawValue is an error string here
			}
		}
		return true, rawValues //All properties were found and were valid
	}

	return false, fmt.Sprintf("Component Type: %v does not exist", componentType)
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
func isType[T interface{}](value interface{}) (bool, interface{}) {
	_, ok := value.(T)
	if !ok {
		return ok, fmt.Sprintf("%v is of type %T; expecting type: %T", value, value, *new(T))
	}
	return ok, value
}

func isInt(value interface{}) (bool, interface{}) {
	//cast to float64
	if float, ok := value.(float64); ok && float == math.Trunc(float) {
		return true, value
	}
	return false, fmt.Sprintf("%v is of type %T; expecting type %T", value, value, *new(int))
}

func isArray(value interface{}) (bool, interface{}) {
	if array, ok := value.([]interface{}); ok {
		var rawArray = make([]interface{}, 0) //initlizing like this so empty will marshall to `[]` not `nil`
		//Have to convert each item one by one
		for _, genericItem := range array {
			item := *new(ArrayItem)
			if error := mapstructure.Decode(genericItem, &item); error != nil {
				return false, fmt.Sprintf("%v", error)
			} else if item.ComponentType == nil || item.Value == nil {
				return false, fmt.Sprintf("%v has a nil value", item)
			} else {
				valid, result := jsonIsValid(*item.ComponentType, *item.Value)
				if valid {
					rawArray = append(rawArray, result)
				} else {
					return valid, result
				}
			}
		}
		return true, rawArray
	} else {
		return false, fmt.Sprintf("%v is of type %T; expecting type %T", value, value, *new([]interface{}))
	}
}

// the key is the name of the component, value is a dictionary of its properties
//
//	the key is the name: the value is the type
var customComponents = make(map[string]map[string]string)

var primitiveComponents map[string]func(interface{}) (bool, interface{})

// The user's JSON
// key: the endpoint path like "/config" or "homepage/version3"
// value: File to JSOn struct
// Future Improvements: the key is repeated in the FileToJson struct "path", make this a Set based on `Path` property uniquness
var files = map[string]FileToJson{} //starts empty
