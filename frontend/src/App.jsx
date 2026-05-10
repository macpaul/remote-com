import {useState, useEffect} from 'react';
import './App.css';
import {ListPorts, GetBindings, AddBinding, RemoveBinding, StartBinding, StopBinding} from "../wailsjs/go/main/App";

function App() {
    const [ports, setPorts] = useState([]);
    const [bindings, setBindings] = useState({});
    const [newBinding, setNewBinding] = useState({serialPort: '', tcpPort: '', password: ''});
    const [error, setError] = useState('');

    const refreshData = () => {
        ListPorts().then(setPorts).catch(err => setError("Failed to list ports: " + err));
        GetBindings().then(setBindings).catch(err => setError("Failed to get bindings: " + err));
    };

    useEffect(() => {
        refreshData();
        const interval = setInterval(refreshData, 5000);
        return () => clearInterval(interval);
    }, []);

    const handleAdd = () => {
        if (!newBinding.serialPort || !newBinding.tcpPort || !newBinding.password) {
            setError("All fields are required");
            return;
        }
        AddBinding(newBinding.serialPort, parseInt(newBinding.tcpPort), newBinding.password)
            .then(() => {
                setNewBinding({serialPort: '', tcpPort: '', password: ''});
                refreshData();
                setError('');
            })
            .catch(err => setError("Failed to add binding: " + err));
    };

    const handleAction = (action, key) => {
        action(key).then(refreshData).catch(err => setError("Action failed: " + err));
    };

    return (
        <div id="App">
            <header className="header">
                <h1>Remote-COM</h1>
            </header>

            <main className="main-content">
                {error && <div className="error-banner">{error}</div>}

                <section className="section">
                    <h2>Add New Binding</h2>
                    <div className="form">
                        <select 
                            value={newBinding.serialPort} 
                            onChange={e => setNewBinding({...newBinding, serialPort: e.target.value})}
                        >
                            <option value="">Select Serial Port</option>
                            {ports.map(p => <option key={p.name} value={p.name}>{p.name}</option>)}
                        </select>
                        <input 
                            type="number" 
                            placeholder="TCP Port" 
                            value={newBinding.tcpPort}
                            onChange={e => setNewBinding({...newBinding, tcpPort: e.target.value})}
                        />
                        <input 
                            type="password" 
                            placeholder="SSH Password" 
                            value={newBinding.password}
                            onChange={e => setNewBinding({...newBinding, password: e.target.value})}
                        />
                        <button onClick={handleAdd}>Add</button>
                    </div>
                </section>

                <section className="section">
                    <h2>Active Bindings</h2>
                    <table className="bindings-table">
                        <thead>
                            <tr>
                                <th>Serial Port</th>
                                <th>TCP Port</th>
                                <th>Status</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {Object.entries(bindings).map(([key, b]) => (
                                <tr key={key}>
                                    <td>{b.serialPort}</td>
                                    <td>{b.tcpPort}</td>
                                    <td>
                                        <span className={`status-pill ${b.active ? 'active' : 'inactive'}`}>
                                            {b.active ? 'Running' : 'Stopped'}
                                        </span>
                                    </td>
                                    <td>
                                        {!b.active ? (
                                            <button onClick={() => handleAction(StartBinding, key)}>Start</button>
                                        ) : (
                                            <button onClick={() => handleAction(StopBinding, key)}>Stop</button>
                                        )}
                                        <button className="btn-danger" onClick={() => handleAction(RemoveBinding, key)}>Remove</button>
                                    </td>
                                </tr>
                            ))}
                            {Object.keys(bindings).length === 0 && (
                                <tr>
                                    <td colSpan="4">No bindings configured.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </section>
            </main>
        </div>
    );
}

export default App;
