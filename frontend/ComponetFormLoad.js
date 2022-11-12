import React from "react";
import { ComponentForm } from "./ComponentForm";

export class ComponentFormLoad extends React.Component {
    constructor(props) {
        super(props);
        this.state = {  
            state: "completed", //"loading", "completed", "error"
            components: ["INT", "STRING", "BOOL", "FLOAT"], //when state is "completed" this will be set with values
        };
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

