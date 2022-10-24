import React from "react";
import { PropertyDefinition } from "./PropertyDefinition";

export class DropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            properties: [
                {
                    name: "property name",
                    type: this.types[0],
                },
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
    }

    types = ["INT", "STRING", "BOOL", "FLOAT"]


    handleNameChange(index, event) {
        this.state.properties[index].name = event.target.value

        this.setState(this.state)
    }

    handleTypeChange(index, event) {
        this.state.properties[index].type = event.target.value
        this.setState(this.state);
    }

    handleSubmit(event) {
        const message = "Want a component with" + this.state.properties.map ( property =>
            "\n a property of type " + property.type + " called " + property.name + ";"
        ).join()

        alert(message);
        event.preventDefault();
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

    render() {
        return (
        <div>
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

