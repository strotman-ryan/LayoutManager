# LayoutManager

Create a server and front end that can create, edit, and serve JSON files


API design 
- Will use HTTP 


- Create component
    - POST 
    - /createcomponent
    - data values
        - name
            - name of the component
            - if this name already exists, replace that compoenent with this component
                - will be used to edit components
        - set of tuples of type (string, string)
            - first string is the JSON object
            - second string is the type of object to be made
                - must equal "string", "bool", "array" etc. or any user made components that have already been made
    - example
        - name: "Address"
        - values: {("street", "string"), ("areaCode", "int"), ("city", "string")}
    - example 
        - name: "Employee"
        - values: {("name", "string"), ("address", "Address"), ("cars", "array")}


- Create file
    - POST
    - /createfile
    - data values
        - name
            - name of the file
            - eventually the endpoint to hit to get that file
            - if name already exists replace it. this will be editing as well
        - contents
            - the raw json to save and serve when hit
        - structure
            - a JSON structure that details what components were used
            - used to verify
    - rules
        - if the contents do not match the given structure -> fail
    - example
        - name: "Employees"
        - contents:
            [
                {
                    "name": "mark smith",
                    "address": {
                        "street": "123 seaseme street",
                        "areaCode": 12345,
                        "city": "cincinnti"
                    },
                    "cars": ["ford","toyota"]
                },
                {
                    "name": "carmen andrews",
                    "address": {
                        "street": "145 look street",
                        "areaCode": 513,
                        "city": "los angles"
                    },
                    "cars": []
                },
                {
                    "name": "cool guy",
                    "address": {
                        "street": "45 yaya",
                        "areaCode": 53253,
                        "city": "Narobi"
                    },
                    "cars": ["Honda"]
                }
            ]
        - structure:
            {
                "type": "Array",
                "contents": [
                {
                    "type": "Employee"
                    "contents": {
                        "name": "string",
                        "address": {
                            "type": "Address",
                            "contents": {
                                "street": "string",
                                "areaCode": "Int",
                                "city": "string"
                            }
                        },
                        "cars": {
                            "type": "Array",
                            "contents": [
                                "string",
                                "string"
                            ] 
                        }
                    }
                },
                {
                    "type": "Employee"
                    "contents": {
                        "name": "string",
                        "address": {
                            "type": "Address",
                            "contents": {
                                "street": "string",
                                "areaCode": "Int",
                                "city": "string"
                            }
                        },
                        "cars": {
                            "type": "Array",
                            "contents": [] 
                        }
                    }
                },
                {
                    "type": "Employee"
                    "contents": {
                        "name": "string",
                        "address": {
                            "type": "Address",
                            "contents": {
                                "street": "string",
                                "areaCode": "Int",
                                "city": "string"
                            }
                        },
                        "cars": {
                            "type": "Array",
                            "contents": [
                                "string"
                            ] 
                        }
                    }
                }]
            }
    - example2
        - contents:
            45
        - structure:
            {
                "type": "Int"
            }
    - example3
        - contents:
            [ 45, "hello"]
        - structure:
            {
                "type": "Array",
                "contents": [
                    "Int",
                    "string"
                ]
            }
    
    
