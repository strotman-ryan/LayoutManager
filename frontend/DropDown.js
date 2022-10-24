import React from "react";
import { PropertyDefinition } from "./PropertyDefinition";

export class DropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            properties: [
                {
                    propertyName: "property name",
                    type: this.types[0],
                },
                {
                    propertyName: "property name",
                    type: this.types[0],
                }
            ]        
        };

        this.handleTypeChange = this.handleTypeChange.bind(this);
        this.handleNameChange = this.handleNameChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    types = ["INT", "STRING", "BOOL", "FLOAT"]


    handleNameChange(index, event) {
        this.state.properties[index].propertyName = event.target.value

        this.setState(this.state)
    }

    handleTypeChange(index, event) {
        this.state.properties[index].type = event.target.value
        this.setState(this.state);
    }

    handleSubmit(event) {
        const message = "Want a component with" + this.state.properties.map ( property =>
            "\n a property of type " + property.type + " called " + property.propertyName + ";"
        ).join()

        alert(message);
        event.preventDefault();
    }

    render() {
        return (
        <form onSubmit={this.handleSubmit}>
            {
                this.state.properties.map ( (property, index) =>
                <PropertyDefinition 
                    type = {property.type}
                    types = {this.types}
                    textValue = {property.propertyName}
                    onTypeChange = {(e) => this.handleTypeChange(index, e)} 
                    onNameChange = {(e) => this.handleNameChange(index, e)}
                />
                )
            }
            
            <input type="submit" value="Submit" />
        </form>
        );
    }
}

