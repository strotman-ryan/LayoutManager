from flask import Flask, request, abort
import json

app = Flask(__name__)

components = {}
supportedJsonPrimitives = {"STRING", "INT", "BOOLEAN"}

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"



@app.route("/component", methods=["POST"])
def component():
    if request.method == "POST":
        name = request.form["name"]
        value = request.form["value"]
        jsonValue = json.loads(value)
        for _, value in jsonValue.items():
            print(value)
            supportedValueTypes = supportedJsonPrimitives.union(components.keys())
            if value not in supportedValueTypes:
                abort(404)
        print(jsonValue)
        components[name] = jsonValue
    return components
