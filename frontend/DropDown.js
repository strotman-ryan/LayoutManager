import React from "react";

export class DropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {value: this.options[0]};

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    options = ["INT", "STRING", "BOOL", "FLOAT"]

    handleChange(event) {
        this.setState({value: event.target.value});
    }

    handleSubmit(event) {
        alert('Your favorite flavor is: ' + this.state.value);
        event.preventDefault();
    }

    render() {
        return (
        <form onSubmit={this.handleSubmit}>
            <label>
            Pick your favorite flavor:
            <PropertyDefinition 
                defaultValue = {this.state.value}
                options = {this.options}
                handleChange = {this.handleChange}  />
            </label>
            <input type="submit" value="Submit" />
        </form>
        );
    }
}

function PropertyDefinition(props) {
    return (
        <div>
            <select value={props.defaultValue} onChange={props.handleChange}>
            {
                props.options.map (option => 
                    <option value={option}> {option} </option>
                )
            }
            </select>
            <input type="text" value="Some text..."/>
        </div>
    );
}