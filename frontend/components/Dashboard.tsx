"use client";

import { useEffect, useState } from "react";

interface Risk {
  id: string;
  title: string;
  score: number;
  status: string;
}

interface Evidence {
  framework: string;
  timestamp: string;
}

export default function Dashboard() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [risks, setRisks] = useState<Risk[]>([]);
  const [evidence, setEvidence] = useState<Evidence | null>(null);

  useEffect(() => {
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "risk_update") {
        setRisks(Object.values(msg.risks));
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  const fetchEvidence = async (format: string) => {
    const res = await fetch(`/api/v1/compliance/evidence?format=${format}`);
    if (format === "json") {
      const data = await res.json();
      setEvidence(data);
    } else {
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `evidence.${format}`;
      a.click();
    }
  };

  const createRisk = async () => {
    await fetch("/api/v1/risks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        title: "Demo Risk",
        description: "Demo description",
        likelihood: "possible",
        impact: "moderate",
      }),
    });
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">Phoenix Dashboard</h1>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Risk Register</h2>
        <button
          onClick={createRisk}
          className="mb-4 px-4 py-2 bg-phoenix-orange text-white rounded"
        >
          Add Risk
        </button>
        <ul className="space-y-2">
          {risks.map((r) => (
            <li key={r.id} className="border p-2 rounded">
              <strong>{r.title}</strong> (Score: {r.score}) – {r.status}
            </li>
          ))}
        </ul>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Compliance Evidence</h2>
        <div className="space-x-4">
          <button
            onClick={() => fetchEvidence("json")}
            className="px-4 py-2 bg-blue-600 text-white rounded"
          >
            Export JSON
          </button>
          <button
            onClick={() => fetchEvidence("pdf")}
            className="px-4 py-2 bg-green-600 text-white rounded"
          >
            Export PDF
          </button>
        </div>
        {evidence && (
          <pre className="mt-4 p-4 bg-gray-100 rounded text-sm">
            {JSON.stringify(evidence, null, 2)}
          </pre>
        )}
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-4">Chat Bot Commands</h2>
        <p className="text-sm text-gray-600">
          Slack: /phoenix-status, /phoenix-risk, /phoenix-alert
        </p>
      </section>
    </main>
  );
}
