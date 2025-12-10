import { useEffect, useRef, useState } from 'react';
import { ClipboardSetText } from "../../wailsjs/runtime";

interface ConsoleProps {
    statusHistory: string[];
    logs: string[];
    onClose: () => void;
}

export function Console({ statusHistory, logs, onClose }: ConsoleProps) {
    const endRef = useRef<HTMLDivElement>(null);
    const [autoScroll, setAutoScroll] = useState(true);

    useEffect(() => {
        if (autoScroll) {
            endRef.current?.scrollIntoView({ behavior: "smooth" });
        }
    }, [logs, statusHistory, autoScroll]);

    const copyToClipboard = () => {
        let allLogs = statusHistory.join("\n");
        if (statusHistory.length > 0 && logs.length > 0) {
            allLogs += "\n--- GAME OUTPUT ---\n";
        }
        allLogs += logs.join("");
        ClipboardSetText(allLogs);
    };

    return (
        <div className="console-overlay">
            <div className="console-window">
                <div className="console-header">
                    <span>GAME LOGS</span>
                    <div className="console-actions">
                        <label className="checkbox-label" style={{ marginRight: '15px', color: '#ccc', fontSize: '0.9em' }}>
                            <input
                                type="checkbox"
                                checked={autoScroll}
                                onChange={(e) => setAutoScroll(e.target.checked)}
                            />
                            Autoscroll
                        </label>
                        <button onClick={copyToClipboard} className="action-btn">COPY</button>
                        <button onClick={onClose} className="close-btn">×</button>
                    </div>
                </div>
                <div className="console-content">
                    {statusHistory.map((line, i) => (
                        <div key={`status-${i}`} className="log-line launcher-status">{line}</div>
                    ))}
                    {statusHistory.length > 0 && logs.length > 0 && (
                        <div className="log-separator">--- GAME OUTPUT ---</div>
                    )}
                    {logs.map((line, i) => (
                        <div key={`log-${i}`} className="log-line">{line}</div>
                    ))}
                    <div ref={endRef} />
                </div>
            </div>
        </div>
    );
}
