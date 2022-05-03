from flask import Flask, request

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"



@app.route("/component", methods=["POST"])
def component():
    if request.method == "POST":
        componentName = request.form["name"]
    return componentName
