import React from "react";
import { PropertyDefinition } from "./PropertyDefinition";
import { backendBaseUrl } from "./App";

export class ComponentForm extends React.Component {
    constructor(props) {
        super(props);
        this.types = props.components
        this.state = {
            title: "Component Name...",
            properties: [
                {
                    name: "property name",
                    type: this.types[0],
                }
            ]        
        };

        this.handleTypeChange = this.handleTypeChange.bind(this);
        this.handleNameChange = this.handleNameChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onAddComponent = this.onAddComponent.bind(this);
        this.onDeleteComponent = this.onDeleteComponent.bind(this);
        this.onTitleChange = this.onTitleChange.bind(this);
    }


    handleNameChange(index, event) {
        this.state.properties[index].name = event.target.value

        this.setState(this.state)
    }

    handleTypeChange(index, event) {
        this.state.properties[index].type = event.target.value
        this.setState(this.state);
    }

    handleSubmit(event) {
        //creating these pass through variables b/c compiler yells at me 
        const title = this.state.title
        var properties = {}
        for (let property in this.state.properties) {
            properties[property.name] = property.type
        }

        var jsonObj = {}
        jsonObj[title] = properties

        fetch(
            backendBaseUrl + "component",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(jsonObj)
            }
        ).then((res) => res.text()) 
        .then((text) => {
            alert(text);
        })
    }

    onAddComponent(event) {
        this.state.properties.push (
            {
                name: "property name",
                type: this.types[0],
            }
        )
        this.setState(this.state)
    }

    onDeleteComponent(index, event) {
        this.state.properties.splice(index, 1)
        this.setState(this.state)
    }

    onTitleChange(event) {
        this.setState({
            "title": event.target.value
        })
    }

    render() {
        return (
        <div>
            <input type="text" value={this.state.title} onChange={this.onTitleChange}/>
            {
                this.state.properties.map ( (property, index) =>
                <PropertyDefinition 
                    type = {property.type}
                    types = {this.types}
                    textValue = {property.name}
                    onDelete = {(e) => this.onDeleteComponent(index, e)}
                    onTypeChange = {(e) => this.handleTypeChange(index, e)} 
                    onNameChange = {(e) => this.handleNameChange(index, e)}
                />
                )
            }
            <button onClick={this.onAddComponent}> Add Property </button>
            <button onClick={this.handleSubmit}> Submit </button>
        </div>

        );
    }
}

