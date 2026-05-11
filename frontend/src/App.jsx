import {useState, useEffect} from 'react';
import './App.css';
import {ListPorts, GetBindings, AddBinding, RemoveBinding, StartBinding, StopBinding, SaveConfig} from "../wailsjs/go/main/App";

function App() {
    const [ports, setPorts] = useState([]);
    const [bindings, setBindings] = useState({});
    const [newBinding, setNewBinding] = useState({
        serialPort: '', 
        tcpPort: '', 
        password: '',
        baudRate: 9600,
        dataBits: 8,
        parity: 'none',
        stopBits: '1',
        flowControl: 'none',
        charDelay: 0,
        lineDelay: 0
    });
    const [error, setError] = useState('');
    const [saveStatus, setSaveStatus] = useState('');

    const refreshData = () => {
        ListPorts().then(setPorts).catch(err => setError("Failed to list ports: " + err));
        GetBindings().then(setBindings).catch(err => setError("Failed to get bindings: " + err));
    };

    useEffect(() => {
        refreshData();
        const interval = setInterval(refreshData, 5000);
        return () => clearInterval(interval);
    }, []);

    const handleSave = () => {
        SaveConfig()
            .then(() => {
                setSaveStatus('Settings saved to settings.ini');
                setTimeout(() => setSaveStatus(''), 3000);
                setError('');
            })
            .catch(err => setError("Failed to save settings: " + err));
    };

    const handleAdd = () => {
        if (!newBinding.serialPort || !newBinding.tcpPort || !newBinding.password) {
            setError("All fields are required");
            return;
        }
        const port = parseInt(newBinding.tcpPort);
        if (isNaN(port) || port < 1 || port > 65535) {
            setError("TCP Port must be between 1 and 65535");
            return;
        }

        const serialConf = {
            baudRate: parseInt(newBinding.baudRate),
            dataBits: parseInt(newBinding.dataBits),
            parity: newBinding.parity,
            stopBits: newBinding.stopBits,
            flowControl: newBinding.flowControl,
            charDelay: parseInt(newBinding.charDelay) || 0,
            lineDelay: parseInt(newBinding.lineDelay) || 0
        };

        AddBinding(newBinding.serialPort, port, newBinding.password, serialConf)
            .then(() => {
                setNewBinding({
                    serialPort: '', 
                    tcpPort: '', 
                    password: '',
                    baudRate: 9600,
                    dataBits: 8,
                    parity: 'none',
                    stopBits: '1',
                    flowControl: 'none',
                    charDelay: 0,
                    lineDelay: 0
                });
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
                <div className="header-actions">
                    {saveStatus && <span className="save-status">{saveStatus}</span>}
                    <button className="btn-save" onClick={handleSave}>Save Settings</button>
                </div>
            </header>

            <main className="main-content">
                {error && <div className="error-banner">{error}</div>}

                <section className="section">
                    <h2>Add New Binding</h2>
                    <div className="form">
                        <div className="form-row">
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
                                min="1"
                                max="65535"
                                value={newBinding.tcpPort}
                                onChange={e => setNewBinding({...newBinding, tcpPort: e.target.value})}
                            />
                            <input 
                                type="password" 
                                placeholder="SSH Password" 
                                value={newBinding.password}
                                onChange={e => setNewBinding({...newBinding, password: e.target.value})}
                            />
                        </div>
                        <div className="form-row">
                            <label>Speed:
                                <select value={newBinding.baudRate} onChange={e => setNewBinding({...newBinding, baudRate: e.target.value})}>
                                    {[110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 57600, 115200, 230400, 460800, 921600].map(b => (
                                        <option key={b} value={b}>{b}</option>
                                    ))}
                                </select>
                            </label>
                            <label>Data:
                                <select value={newBinding.dataBits} onChange={e => setNewBinding({...newBinding, dataBits: e.target.value})}>
                                    <option value="7">7bit</option>
                                    <option value="8">8bit</option>
                                </select>
                            </label>
                            <label>Parity:
                                <select value={newBinding.parity} onChange={e => setNewBinding({...newBinding, parity: e.target.value})}>
                                    <option value="none">none</option>
                                    <option value="odd">odd</option>
                                    <option value="even">even</option>
                                    <option value="mark">mark</option>
                                    <option value="space">space</option>
                                </select>
                            </label>
                            <label>Stop bits:
                                <select value={newBinding.stopBits} onChange={e => setNewBinding({...newBinding, stopBits: e.target.value})}>
                                    <option value="1">1bit</option>
                                    <option value="1.5">1.5bit</option>
                                    <option value="2">2bit</option>
                                </select>
                            </label>
                        </div>
                        <div className="form-row">
                            <label>Flow control:
                                <select value={newBinding.flowControl} onChange={e => setNewBinding({...newBinding, flowControl: e.target.value})}>
                                    <option value="none">none</option>
                                    <option value="xonxoff">Xon/Xoff</option>
                                    <option value="rtscts">RTS/CTS</option>
                                    <option value="dsrdtr">DSR/DTR</option>
                                </select>
                            </label>
                            <label>Char Delay (ms):
                                <input type="number" min="0" value={newBinding.charDelay} onChange={e => setNewBinding({...newBinding, charDelay: e.target.value})} />
                            </label>
                            <label>Line Delay (ms):
                                <input type="number" min="0" value={newBinding.lineDelay} onChange={e => setNewBinding({...newBinding, lineDelay: e.target.value})} />
                            </label>
                            <button onClick={handleAdd}>Add</button>
                        </div>
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
