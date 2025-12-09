import { useState, useEffect } from 'react';
import './App.css';
import { LaunchGame, GetSystemInfo } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime";
import { Console } from "./components/Console";
import { main } from "../wailsjs/go/models";

function App() {
    const [status, setStatus] = useState("Ready to Launch");
    const [username, setUsername] = useState("Player");
    const [showLogWhileRunning, setShowLogWhileRunning] = useState(false);
    const [logs, setLogs] = useState<string[]>([]);
    const [statusHistory, setStatusHistory] = useState<string[]>([]);
    const [isConsoleOpen, setIsConsoleOpen] = useState(false);
    const [ramMB, setRamMB] = useState(2048);
    const [systemInfo, setSystemInfo] = useState<main.SystemInfo | null>(null);

    useEffect(() => {
        // Fetch system info on startup
        GetSystemInfo().then((info) => {
            setSystemInfo(info);
            setRamMB(info.defaultRAM);
        });

        const unsubscribeStatus = EventsOn("update-status", (msg: string) => {
            setStatus(msg);
            setStatusHistory(prev => [...prev, `[LAUNCHER] ${msg}`]);
            if (msg === "Crashed") {
                setIsConsoleOpen(true);
            } else if (msg === "Ready to Launch") {
                // Auto-hide log when game quits normally (not when crashed)
                setIsConsoleOpen(false);
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
        setStatusHistory([]); // Clear status history on new launch
        if (showLogWhileRunning) {
            setIsConsoleOpen(true);
        }
        LaunchGame(username, ramMB).then(() => {
            // Status updates will come via events
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

                <div className="input-group">
                    <label>
                        RAM ALLOCATION (GiB)
                        {systemInfo && (
                            <span className="ram-info">
                                {systemInfo.is32Bit && " (32-bit limited to 1 GiB)"}
                                {!systemInfo.is32Bit && ` (System: ${Math.floor(systemInfo.totalRAM / 1024)} GiB)`}
                            </span>
                        )}
                    </label>
                    <input
                        type="number"
                        value={Math.round(ramMB / 1024)}
                        onChange={(e) => {
                            const valueGiB = parseInt(e.target.value) || 0;
                            const valueMiB = valueGiB * 1024;
                            if (systemInfo) {
                                const clamped = Math.min(Math.max(valueMiB, systemInfo.minRAM), systemInfo.maxRAM);
                                setRamMB(clamped);
                            } else {
                                setRamMB(valueMiB);
                            }
                        }}
                        min={systemInfo ? Math.ceil(systemInfo.minRAM / 1024) : 1}
                        max={systemInfo ? Math.floor(systemInfo.maxRAM / 1024) : 32}
                        disabled={systemInfo?.is32Bit}
                        className="ram-input"
                        placeholder="2"
                    />
                </div>

                <div className="actions">
                    <div className="options">
                        <label className="checkbox-label">
                            <input
                                type="checkbox"
                                checked={showLogWhileRunning}
                                onChange={(e) => setShowLogWhileRunning(e.target.checked)}
                            />
                            Show Log While Running
                        </label>
                        {(logs.length > 0 || statusHistory.length > 0) && !isConsoleOpen && (
                            <button
                                className="btn-show-log"
                                onClick={() => setIsConsoleOpen(true)}
                            >
                                SHOW LOG
                            </button>
                        )}
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
                    statusHistory={statusHistory}
                    logs={logs}
                    onClose={() => setIsConsoleOpen(false)}
                />
            )}
        </div>
    )
}

export default App
