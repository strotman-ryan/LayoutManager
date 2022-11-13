import React from "react";
import { ComponentFormLoad } from "./ComponetFormLoad";

export function App() {
  return (
    <div>
      <ComponentFormLoad />
      {/* for numbers */}
      <input placeholder="numbers only" type="number"></input>
      {/* for boolean */}
      <input type="checkbox" value="flase"></input>
      {/* for string */}
      <input ></input>
    </div> 
  );
}

//constants
export let backendBaseUrl: String = "http://localhost:8080/"
