import { useState, useEffect } from 'react';
import './App.css';
import { LaunchGame } from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from "../wailsjs/runtime";
import { Console } from "./components/Console";

function App() {
    const [status, setStatus] = useState("Ready to Launch");
    const [username, setUsername] = useState("Player");
    const [showLog, setShowLog] = useState(false);
    const [logs, setLogs] = useState<string[]>([]);
    const [isConsoleOpen, setIsConsoleOpen] = useState(false);

    useEffect(() => {
        const unsubscribeStatus = EventsOn("update-status", (msg: string) => {
            setStatus(msg);
            if (msg === "Crashed") {
                setIsConsoleOpen(true);
            }
        });

        const unsubscribeLogs = EventsOn("log-data", (msg: string) => {
            setLogs(prev => [...prev, msg]);
        });

        return () => {
            unsubscribeStatus();
            unsubscribeLogs();
        };
    }, []);

    const launch = () => {
        setStatus("Launching...");
        setLogs([]); // Clear logs on new launch
        if (showLog) {
            setIsConsoleOpen(true);
        }
        LaunchGame(username).then((result) => {
            // result is "Launching..." usually, status updates will come via events
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
                    <div className="options">
                        <label className="checkbox-label">
                            <input
                                type="checkbox"
                                checked={showLog}
                                onChange={(e) => setShowLog(e.target.checked)}
                            />
                            Show Log
                        </label>
                    </div>

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

            {isConsoleOpen && (
                <Console
                    logs={logs}
                    onClose={() => setIsConsoleOpen(false)}
                />
            )}
        </div>
    )
}

export default App
