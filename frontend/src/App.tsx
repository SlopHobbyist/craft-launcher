import { useState, useEffect } from 'react';
import './App.css';
import { LaunchGame } from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from "../wailsjs/runtime";

function App() {
    const [status, setStatus] = useState("Ready to Launch");
    const [username, setUsername] = useState("Player");

    useEffect(() => {
        const unsubscribe = EventsOn("update-status", (msg: string) => {
            setStatus(msg);
        });
        return () => {
            // EventsOn returns a cleanup function, or we can use EventsOff if we knew the handler ref, 
            // but wailsjs runtime.d.ts says EventsOn returns () => void.
            unsubscribe();
        };
    }, []);

    const launch = () => {
        setStatus("Launching...");
        LaunchGame(username).then((result) => {
            // result is "Launching..." usually, status updates will come via events
            // setStatus(result); // Don't override event status
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
