

export function PropertyDefinition(props) {
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