import { useState } from 'react';
import './App.css';
import { LaunchGame } from "../wailsjs/go/main/App";

function App() {
    const [status, setStatus] = useState("Ready to Launch");
    const [username, setUsername] = useState("Player");

    const launch = () => {
        setStatus("Launching...");
        LaunchGame(username).then((result) => {
            setStatus(result);
        });
    };

    return (
        <div id="App">
            <div className="container">
                <h1 className="title">MINECRAFT 1.8.9</h1>

                <div className="input-group">
                    <label>USERNAME</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        placeholder="Offline Username"
                        className="username-input"
                    />
                </div>

                <div className="actions">
                    <button className="btn secondary" disabled>
                        CHECK FOR UPDATES
                    </button>
                    <button className="btn primary" onClick={launch}>
                        PLAY ({username})
                    </button>
                </div>

                <div className="status-bar">
                    STATUS: {status}
                </div>
            </div>
        </div>
    )
}

export default App
