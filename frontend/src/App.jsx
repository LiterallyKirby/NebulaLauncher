import "./App.css"; // <-- import the CSS file
import React, { useEffect, useState } from "react";
import { GetProcesses, Inject } from "../wailsjs/go/main/App";

function App() {
	const [processes, setProcesses] = useState([]);
	const [log, setLog] = useState([]);

	useEffect(() => {
		GetProcesses().then(setProcesses);
	}, []);

	const handleInject = async (pid) => {
		try {
			await Inject(pid);
			setLog((prev) => [`Injected into PID ${pid}`, ...prev]);
		} catch (err) {
			setLog((prev) => [`Failed to inject PID ${pid}: ${err}`, ...prev]);
		}
	};

	return (
		<div className="Screen">
			<div className="Title">
				Nebula Launcher
				{[...Array(12)].map((_, i) => (
					<div key={i} className={`star star${i + 1}`}></div>
				))}
			</div>

			{/* Scrollable container for buttons */}
			<div className="scroll-container">
				<div className="grid-container">
					{processes.map((p) => (
						<button key={p.pid} onClick={() => handleInject(p.pid)}>
							{p.pid} - {p.name}
						</button>
					))}
				</div>
			</div>

			{/* Output / log section */}
			<div className="log-container">
				{log.map((entry, i) => (
					<div key={i} className="log-entry">
						{entry}
					</div>
				))}
			</div>
		</div>
	);
}

export default App;
