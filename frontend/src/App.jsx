import "./App.css";
import React, { useEffect, useState } from "react";
import { GetProcesses, Inject } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

function App() {
  const [processes, setProcesses] = useState([]);
  const [logs, setLogs] = useState([]); // store logs

  useEffect(() => {
    // Fetch processes initially and every 3 seconds
    const fetchProcesses = () => {
      GetProcesses().then(setProcesses);
    };

    fetchProcesses();
    const interval = setInterval(fetchProcesses, 3000);

    // Listen to backend log events
    EventsOn("log", (msg) => {
      setLogs((prev) => [...prev, msg]);
    });

    return () => clearInterval(interval);
  }, []);

  const handleInject = (pid) => {
    Inject(pid);
  };

  return (
    <div className="Screen">
      <div className="Title">
        Nebula Launcher
        {[...Array(12)].map((_, i) => (
          <div key={i} className={`star star${i + 1}`}></div>
        ))}
      </div>

      <div className="scroll-container">
   
<div className="grid-container">
  {processes
    .filter(
      (p) =>
        p.name.toLowerCase().includes("java") ||
        p.name.toLowerCase().includes("minecraft") ||
        p.name.includes("1.8.9")
    )
    .map((p) => (
      <button key={p.pid} onClick={() => handleInject(p.pid)}>
        {p.pid} - {p.name}
      </button>
    ))}
</div>
</div>

      {/* Log / Status Panel */}
      <div className="log-container">
        {logs.map((log, i) => (
          <div key={i} className="log-entry">{log}</div>
        ))}
      </div>
    </div>
  );
}

export default App;
