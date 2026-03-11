"use client";

import { useEffect, useState } from "react";

interface ISMS {
  id: string;
  name: string;
  description: string;
  version: string;
  status: string;
  compliance_score?: number;
  created_at: string;
}

interface Control {
  id: string;
  number: string;
  title: string;
  category: string;
  status: string;
  owner: string;
  last_tested?: string;
  next_review: string;
}

interface Risk {
  id: string;
  title: string;
  risk_level: string;
  likelihood: string;
  impact: string;
  treatment: string;
  status: string;
  owner: string;
}

export default function ISO27001Dashboard() {
  const [isms, setISMS] = useState<ISMS[]>([]);
  const [selectedISMS, setSelectedISMS] = useState<ISMS | null>(null);
  const [controls, setControls] = useState<Control[]>([]);
  const [risks, setRisks] = useState<Risk[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch ISMS list
    fetch("/api/v1/compliance/iso27001/isms")
      .then((res) => res.json())
      .then((data) => setISMS(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "iso27001_update") {
        setISMS(Object.values(msg.isms));
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedISMS) {
      // Fetch compliance score
      fetch(`/api/v1/compliance/iso27001/isms/${selectedISMS.id}/score`)
        .then((res) => res.json())
        .then((data) => {
          setSelectedISMS((prev) => prev ? { ...prev, compliance_score: data.compliance_score } : null);
        });

      // Fetch controls
      fetch("/api/v1/compliance/iso27001/controls")
        .then((res) => res.json())
        .then((data) => setControls(data));

      // Fetch risks
      fetch("/api/v1/compliance/iso27001/risks")
        .then((res) => res.json())
        .then((data) => setRisks(data));
    }
  }, [selectedISMS]);

  const createISMS = async () => {
    await fetch("/api/v1/compliance/iso27001/isms", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Enterprise ISMS",
        description: "ISO 27001 Information Security Management System",
        version: "2024.1",
      }),
    });
    // Refresh list
    const res = await fetch("/api/v1/compliance/iso27001/isms");
    setISMS(await res.json());
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active": return "text-green-600";
      case "draft": return "text-gray-600";
      case "under_review": return "text-yellow-600";
      case "approved": return "text-blue-600";
      case "non_compliant": return "text-red-600";
      default: return "text-gray-600";
    }
  };

  const getRiskLevelColor = (level: string) => {
    switch (level) {
      case "low": return "bg-green-100 text-green-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "critical": return "bg-red-100 text-red-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">ISO 27001 Compliance</h1>

      {/* ISMS Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={createISMS} className="px-4 py-2 bg-blue-600 text-white rounded">
            Create ISMS
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {isms.map((ismsItem) => (
            <div
              key={ismsItem.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedISMS?.id === ismsItem.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedISMS(ismsItem)}
            >
              <h3 className="font-semibold">{ismsItem.name}</h3>
              <p className="text-sm text-gray-600">{ismsItem.description}</p>
              <div className="flex justify-between items-center mt-2">
                <span className={`text-sm font-medium ${getStatusColor(ismsItem.status)}`}>
                  {ismsItem.status}
                </span>
                <span className="text-sm text-gray-500">{ismsItem.version}</span>
              </div>
              {ismsItem.compliance_score !== undefined && (
                <div className="mt-2">
                  <div className="flex justify-between text-sm">
                    <span>Compliance Score</span>
                    <span className="font-medium">{ismsItem.compliance_score.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full"
                      style={{ width: `${ismsItem.compliance_score}%` }}
                    ></div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      {selectedISMS && (
        <>
          {/* Controls Overview */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Controls</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {controls.filter((c) => c.status === "effective").length}
                </div>
                <div className="text-green-600 text-sm">Effective</div>
              </div>
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {controls.filter((c) => c.status === "implemented").length}
                </div>
                <div className="text-blue-600 text-sm">Implemented</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {controls.filter((c) => c.status === "partially_implemented").length}
                </div>
                <div className="text-yellow-600 text-sm">Partial</div>
              </div>
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {controls.filter((c) => c.status === "not_implemented").length}
                </div>
                <div className="text-red-600 text-sm">Not Implemented</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Control</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Category</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Owner</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Next Review</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {controls.slice(0, 10).map((control) => (
                    <tr key={control.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {control.number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {control.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {control.category}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(control.status)}`}>
                          {control.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {control.owner}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(control.next_review).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Risk Management */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Risk Assessment</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {risks.filter((r) => r.risk_level === "critical").length}
                </div>
                <div className="text-red-600 text-sm">Critical</div>
              </div>
              <div className="bg-orange-50 border border-orange-200 rounded p-4">
                <div className="text-orange-800 text-2xl font-bold">
                  {risks.filter((r) => r.risk_level === "high").length}
                </div>
                <div className="text-orange-600 text-sm">High</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {risks.filter((r) => r.risk_level === "medium").length}
                </div>
                <div className="text-yellow-600 text-sm">Medium</div>
              </div>
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {risks.filter((r) => r.risk_level === "low").length}
                </div>
                <div className="text-green-600 text-sm">Low</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Risk</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Risk Level</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Likelihood</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Impact</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Treatment</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Owner</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {risks.slice(0, 10).map((risk) => (
                    <tr key={risk.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {risk.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskLevelColor(risk.risk_level)}`}>
                          {risk.risk_level}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {risk.likelihood}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {risk.impact}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {risk.treatment}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {risk.owner}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>
        </>
      )}
    </main>
  );
}
