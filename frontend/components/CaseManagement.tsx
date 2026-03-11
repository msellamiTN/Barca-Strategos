"use client";

import { useEffect, useState } from "react";

interface Case {
  id: string;
  title: string;
  status: string;
  severity: string;
  assignee_id: string;
  created_at: string;
  sla?: {
    response_due: string;
    resolve_due: string;
  };
}

interface Action {
  id: string;
  case_id: string;
  type: string;
  status: string;
  executed_at: string;
  executed_by: string;
  details?: string;
}

interface Alert {
  id: string;
  timestamp: string;
  title: string;
  source: string;
  severity: string;
  details: any;
}

export default function CaseManagement() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [cases, setCases] = useState<Case[]>([]);
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [selectedCase, setSelectedCase] = useState<Case | null>(null);

  useEffect(() => {
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "case_update") {
        setCases(Object.values(msg.cases));
      }
      if (msg.type === "siem_alert") {
        setAlerts(msg.alerts);
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  const createCase = async () => {
    await fetch("/api/v1/cases", {
      method: "POST",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${localStorage.getItem("token")}` },
      body: JSON.stringify({
        title: "Phishing Incident",
        description: "Suspicious email reported",
        severity: "high",
        assignee_id: "analyst1",
      }),
    });
  };

  const executePlaybook = async (playbookId: string) => {
    await fetch(`/api/v1/playbooks/${playbookId}/execute`, {
      method: "POST",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${localStorage.getItem("token")}` },
      body: JSON.stringify({ context: { case_id: selectedCase?.id } }),
    });
  };

  const ingestAlert = async (source: string) => {
    await fetch(`/api/v1/siem/ingest/${source}`, {
      method: "POST",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${localStorage.getItem("token")}` },
      body: JSON.stringify({
        _id: "alert-1",
        _time: new Date().toISOString(),
        _raw: "Malware detected on host10",
        _severity: "CRITICAL",
        _source: "edr",
      }),
    });
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">Case Management</h1>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Actions</h2>
        <div className="flex gap-4">
          <button onClick={createCase} className="px-4 py-2 bg-blue-600 text-white rounded">
            Create Case
          </button>
          <button onClick={() => ingestAlert("splunk")} className="px-4 py-2 bg-orange-600 text-white rounded">
            Ingest Splunk Alert
          </button>
        </div>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Cases</h2>
        <ul className="space-y-2">
          {cases.map((c) => (
            <li
              key={c.id}
              className={`border p-4 rounded cursor-pointer ${selectedCase?.id === c.id ? "border-blue-500" : ""}`}
              onClick={() => setSelectedCase(c)}
            >
              <strong>{c.title}</strong> ({c.severity}) – {c.status}
              {c.sla && (
                <div className="text-sm text-gray-600 mt-1">
                  Response due: {new Date(c.sla.response_due).toLocaleString()}
                </div>
              )}
            </li>
          ))}
        </ul>
      </section>

      {selectedCase && (
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Playbooks</h2>
          <button
            onClick={() => executePlaybook("phishing-response")}
            className="px-4 py-2 bg-green-600 text-white rounded"
          >
            Run Phishing Response
          </button>
        </section>
      )}

      <section>
        <h2 className="text-xl font-semibold mb-4">SIEM Alerts</h2>
        <ul className="space-y-2">
          {alerts.slice(0, 5).map((a) => (
            <li key={a.id} className="border p-2 rounded">
              <strong>{a.title}</strong> ({a.severity}) – {a.source}
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
