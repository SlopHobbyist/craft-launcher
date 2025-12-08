import { useEffect, useRef } from 'react';
import { ClipboardSetText } from "../../wailsjs/runtime";

interface ConsoleProps {
    logs: string[];
    onClose: () => void;
}

export function Console({ logs, onClose }: ConsoleProps) {
    const endRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        endRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [logs]);

    const copyToClipboard = () => {
        ClipboardSetText(logs.join(""));
    };

    return (
        <div className="console-overlay">
            <div className="console-window">
                <div className="console-header">
                    <span>GAME LOGS</span>
                    <div className="console-actions">
                        <button onClick={copyToClipboard} className="action-btn">COPY</button>
                        <button onClick={onClose} className="close-btn">×</button>
                    </div>
                </div>
                <div className="console-content">
                    {logs.map((line, i) => (
                        <div key={i} className="log-line">{line}</div>
                    ))}
                    <div ref={endRef} />
                </div>
            </div>
        </div>
    );
}
