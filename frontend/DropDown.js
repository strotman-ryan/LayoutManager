import React from "react";

export class DropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            propertyName: "property name",
            type: this.types[0],
        };

        this.handleTypeChange = this.handleTypeChange.bind(this);
        this.handleNameChange = this.handleNameChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    types = ["INT", "STRING", "BOOL", "FLOAT"]


    handleNameChange(event) {
        this.setState({propertyName: event.target.value})
    }

    handleTypeChange(event) {
        this.setState({type: event.target.value});
    }

    handleSubmit(event) {
        alert("Want a property of type " + this.state.type + " called " + this.state.propertyName);
        event.preventDefault();
    }

    render() {
        return (
        <form onSubmit={this.handleSubmit}>
            <label>
            Pick your favorite flavor:
            <PropertyDefinition 
                type = {this.state.type}
                types = {this.types}
                textValue = {this.state.propertyName}
                onTypeChange = {this.handleTypeChange} 
                onNameChange = {this.handleNameChange}
                 />
            </label>
            <input type="submit" value="Submit" />
        </form>
        );
    }
}

function PropertyDefinition(props) {
    return (
        <div>
            <input type="text" value={props.textValue} onChange={props.onNameChange}/>
            <select value={props.type} onChange={props.onTypeChange}>
            {
                props.types.map (type => 
                    <option value={type}> {type} </option>
                )
            }
            </select>
        </div>
    );
}