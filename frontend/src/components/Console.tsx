import { useEffect, useRef } from 'react';

interface ConsoleProps {
    logs: string[];
    onClose: () => void;
}

export function Console({ logs, onClose }: ConsoleProps) {
    const endRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        endRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [logs]);

    return (
        <div className="console-overlay">
            <div className="console-window">
                <div className="console-header">
                    <span>GAME LOGS</span>
                    <button onClick={onClose} className="close-btn">×</button>
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
