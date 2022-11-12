import React from "react";
import { ComponentForm } from "./ComponentForm";

export class ComponentFormLoad extends React.Component {
    constructor(props) {
        super(props);
        this.state = {  
            state: "loading", //"loading", "completed", "error"
            components: [], //when state is "completed" this will be set with values
        };
    }

    componentDidMount() {
        //get all available components
        fetch("http://localhost:8080/component")
        .then((res) => res.json()) 
        .then((components) => {
            //components is [string]
            this.state.state = "completed"
            this.state.components = components
            this.setState(this.state)
        })
    }
            

    render() {
        if (this.state.state == "loading") {
            return <p>loading...</p>
        } 
        if (this.state.state == "completed") {
            return <ComponentForm components={this.state.components}/>
        }
        //some error occured
        return <p>ERROR</p>
    }
}

