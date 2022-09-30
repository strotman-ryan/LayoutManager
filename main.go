package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/component", getCustomComponentsGin)
	router.POST("/component")

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON. prefix with gin
func getCustomComponentsGin(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getCustomComponents())
}

func addCustomComponentGin(c *gin.Context) {

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
		if !slices.Contains(componentNames, value) {
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
