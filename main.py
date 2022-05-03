from flask import Flask, request, abort
import json

app = Flask(__name__)

components = {}

#TODO:  add support for Arrays
supportedJsonPrimitives = {"STRING", "NUMBER", "BOOLEAN", "NULL"}


@app.route("/component/<name>", methods=["GET","POST"])
def component(name):
    if request.method == "POST":
        value = request.form["value"]
        jsonValue = json.loads(value)
        for _, value in jsonValue.items():
            if value not in supportedJsonPrimitives.union(components.keys()):
                abort(404)
        components[name] = jsonValue
        return components
    
    if request.method == "GET":
        return components[name]

@app.route("/allcomponents")
def getAllComponent():
    return components
