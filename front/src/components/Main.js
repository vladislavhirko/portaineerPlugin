import React, {useState, useEffect, useCallback} from 'react'
import axios from 'axios'

function Jopa(){
    const [container, setContainer] = useState("");
    const [chanel, setChanel] = useState("");
    const [containerList, setContainerList] = useState([]);

    const handleItem = (item) => () => {
        alert(JSON.stringify(item, null, 3))
    }

    console.log(containerList)
    useEffect(() => {
        axios.get('http://localhost:1119/containers').then(responce => setContainerList(responce.data))
    }, [])

    const printToConsole = (e) =>{
        e.preventDefault();
    }

    return (
            <div>
                <div>
                    <form>
                        <div className="row">
                            <div className="col">
                                <div className="form-group">
                                    <label htmlFor="docker-container-name">Container name</label>
                                    <input type="text" 
                                    name="docker-container-name" 
                                    className="form-control" 
                                    id="docker-container-name" 
                                    value={container} 
                                    onChange={e => setContainer(e.target.value)}
                                    placeholder="Container name" />
                                </div>
                            </div>
                            <div className="col">
                                <div className="form-group">
                                    <label htmlFor="chanel-name">Chanel name</label>
                                    <input type="text" 
                                    name="chanel-name" 
                                    className="form-control" 
                                    id="chanel-name" 
                                    value={chanel} 
                                    onChange={e => setChanel(e.target.value)}
                                    placeholder="Cnanel name" />
                                </div>
                            </div>
                        </div>
                        <div className="row">
                            <div className="col">
                                <button className="btn btn-light btn-block button-border" onClick={printToConsole}>ADD</button>
                            </div>
                        </div>
                    </form>

                    <div className="key-value-table">
                        <table className="table">
                            <thead>
                            <tr>
                                <th> Container </th>
                                <th> Chanel </th>
                            </tr>    
                            </thead>   
                            <tbody>
                            {containerList.map(((item) => {
                                return (
                                <tr key={item.Id} onClick={handleItem(item)}>
                                    <td>
                                        {item.Names[0]}
                                    </td>
                                    <td>
                                    {item.Image}
                                    </td>
                                </tr>
                                )
                            }))}
                            </tbody>
                        </table>
                    </div>
                </div>    
            </div>      
    )
}

export default Jopa