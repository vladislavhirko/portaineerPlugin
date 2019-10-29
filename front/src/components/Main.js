import React, {useState} from 'react'

function Jopa(){

    const [container, setContainer] = useState("");
    const [chanel, setChanel] = useState("");

    var first = "";
    var second = "";

    const putData = (e) =>{
        e.preventDefault();
        first = container;
        second = chanel;
        console.log(first, second)
    }

    return (
        <html>
            <head>
                <meta charSet="utf-8" />
                <title>Test</title>

                <meta name="viewport" content="width=device-width, initial-scale=1" />

                {/* <link rel="stylesheet" href="index-style.css" /> */}
                <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous" />
                <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
                <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
                <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
            </head>
            <body>
                <h2>First value - {container}</h2>
                <h2>Second value - {chanel}</h2>
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
                                    value={chanel} 
                                    onChange={e => setChanel(e.target.value)}
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
                                    value={container} 
                                    onChange={e => setContainer(e.target.value)}
                                    placeholder="Cnanel name" />
                                </div>
                            </div>
                        </div>
                        <div className="row">
                            <div className="col">
                                <button className="btn btn-light btn-block button-border" onClick={putData}>ADD</button>
                            </div>
                        </div>
                    </form>

                    <div className="key-value-table">
                        <table className="table">
                            <tr>
                                <th> Container </th>
                                <th> Chanel </th>
                            </tr>
                            <tr>
                                <td>
                                    Lorem
                                </td>
                                <td>
                                    Ipsum
                                </td>
                            </tr>
                            <tr>
                                <td>
                                    Lorem
                                </td>
                                <td>
                                    Ipsum
                                </td>
                            </tr>
                        </table>
                    </div>
                </div>
            </body>
        </html>            
    )
}

export default Jopa