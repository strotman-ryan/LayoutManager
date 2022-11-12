import React from "react";
import { ComponentForm } from "./ComponentForm";
import { backendBaseUrl } from "./App";

interface ComponentFormLoadState {
    state: string;
    components: string[];
}

export class ComponentFormLoad extends React.Component<{}, ComponentFormLoadState> {
    constructor(props: {}) {
        super(props);
        this.state = {  
            state: "loading", //"loading", "completed", "error"
            components: [], //when state is "completed" this will be set with values
        };
    }

    componentDidMount() {
        //get all available components
        fetch( backendBaseUrl + "component")
        .then((res) => res.json()) 
        .then((components) => {
            //components is [string]
            this.setState(
                {
                    state: "completed",
                    components: components
                }
            )
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

