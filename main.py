from flask import Flask, request, abort
import json

app = Flask(__name__)

components = {}

#TODO:  add support for Arrays
supportedJsonPrimitives = {"STRING", "NUMBER", "BOOLEAN", "NULL"}


@app.route("/component", methods=["POST"])
def component():
    if request.method == "POST":
        name = request.form["name"]
        value = request.form["value"]
        jsonValue = json.loads(value)
        for _, value in jsonValue.items():
            if value not in supportedJsonPrimitives.union(components.keys()):
                abort(404)
        components[name] = jsonValue
    return "It worked"
