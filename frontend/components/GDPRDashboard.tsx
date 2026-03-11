"use client";

import { useEffect, useState } from "react";

interface GDPR {
  id: string;
  name: string;
  description: string;
  version: string;
  status: string;
  compliance_score?: number;
  created_at: string;
}

interface Principle {
  id: string;
  title: string;
  article: string;
  category: string;
  status: string;
  owner: string;
  next_review: string;
}

interface Process {
  id: string;
  name: string;
  description: string;
  purpose: string;
  legal_basis: string;
  status: string;
  owner: string;
  next_review: string;
}

interface Request {
  id: string;
  type: string;
  data_subject: string;
  status: string;
  priority: string;
  assigned_to: string;
  due_date: string;
  created_at: string;
}

interface Breach {
  id: string;
  title: string;
  type: string;
  date: string;
  detected: string;
  affected: number;
  status: string;
  notification_required: boolean;
}

interface DPIA {
  id: string;
  title: string;
  description: string;
  controller: string;
  assessor: string;
  status: string;
  review_date: string;
}

export default function GDPRDashboard() {
  const [gdpr, setGDPR] = useState<GDPR[]>([]);
  const [selectedGDPR, setSelectedGDPR] = useState<GDPR | null>(null);
  const [principles, setPrinciples] = useState<Principle[]>([]);
  const [processes, setProcesses] = useState<Process[]>([]);
  const [requests, setRequests] = useState<Request[]>([]);
  const [breaches, setBreaches] = useState<Breach[]>([]);
  const [dpia, setDPIA] = useState<DPIA[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch GDPR list
    fetch("/api/v1/compliance/gdpr")
      .then((res) => res.json())
      .then((data) => setGDPR(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "gdpr_update") {
        setGDPR(Object.values(msg.gdpr));
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedGDPR) {
      // Fetch compliance score
      fetch(`/api/v1/compliance/gdpr/${selectedGDPR.id}/score`)
        .then((res) => res.json())
        .then((data) => {
          setSelectedGDPR((prev) => prev ? { ...prev, compliance_score: data.compliance_score } : null);
        });

      // Fetch processes
      fetch("/api/v1/compliance/gdpr/processes")
        .then((res) => res.json())
        .then((data) => setProcesses(data));

      // Fetch requests
      fetch("/api/v1/compliance/gdpr/requests")
        .then((res) => res.json())
        .then((data) => setRequests(data));

      // Fetch breaches
      fetch("/api/v1/compliance/gdpr/breaches")
        .then((res) => res.json())
        .then((data) => setBreaches(data));

      // Fetch DPIA
      fetch("/api/v1/compliance/gdpr/dpias")
        .then((res) => res.json())
        .then((data) => setDPIA(data));
    }
  }, [selectedGDPR]);

  const createGDPR = async () => {
    await fetch("/api/v1/compliance/gdpr", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Enterprise GDPR Compliance",
        description: "GDPR compliance management system",
        version: "2024.1",
      }),
    });
    // Refresh list
    const res = await fetch("/api/v1/compliance/gdpr");
    setGDPR(await res.json());
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

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "low": return "bg-green-100 text-green-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "critical": return "bg-red-100 text-red-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const getPrincipleColor = (category: string) => {
    switch (category) {
      case "lawfulness": return "bg-blue-50 border-blue-200";
      case "fairness": return "bg-green-50 border-green-200";
      case "transparency": return "bg-purple-50 border-purple-200";
      case "purpose": return "bg-orange-50 border-orange-200";
      case "data_minimization": return "bg-red-50 border-red-200";
      case "accuracy": return "bg-indigo-50 border-indigo-200";
      case "storage_limit": return "bg-yellow-50 border-yellow-200";
      case "security": return "bg-pink-50 border-pink-200";
      case "accountability": return "bg-gray-50 border-gray-200";
      default: return "bg-gray-50 border-gray-200";
    }
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">GDPR Compliance</h1>

      {/* GDPR Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={createGDPR} className="px-4 py-2 bg-blue-600 text-white rounded">
            Create GDPR
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {gdpr.map((gdprItem) => (
            <div
              key={gdprItem.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedGDPR?.id === gdprItem.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedGDPR(gdprItem)}
            >
              <h3 className="font-semibold">{gdprItem.name}</h3>
              <p className="text-sm text-gray-600">{gdprItem.description}</p>
              <div className="flex justify-between items-center mt-2">
                <span className={`text-sm font-medium ${getStatusColor(gdprItem.status)}`}>
                  {gdprItem.status}
                </span>
                <span className="text-sm text-gray-500">{gdprItem.version}</span>
              </div>
              {gdprItem.compliance_score !== undefined && (
                <div className="mt-2">
                  <div className="flex justify-between text-sm">
                    <span>Compliance Score</span>
                    <span className="font-medium">{gdprItem.compliance_score.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-green-600 h-2 rounded-full"
                      style={{ width: `${gdprItem.compliance_score}%` }}
                    ></div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      {selectedGDPR && (
        <>
          {/* Principles Overview */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Data Protection Principles</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {principles.filter((p) => p.status === "compliant").length}
                </div>
                <div className="text-green-600 text-sm">Compliant</div>
              </div>
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {principles.filter((p) => p.status === "implemented").length}
                </div>
                <div className="text-blue-600 text-sm">Implemented</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {principles.filter((p) => p.status === "partially_implemented").length}
                </div>
                <div className="text-yellow-600 text-sm">Partial</div>
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {principles.slice(0, 6).map((principle) => (
                <div key={principle.id} className={`border rounded-lg p-4 ${getPrincipleColor(principle.category)}`}>
                  <h3 className="font-semibold text-sm">{principle.title}</h3>
                  <p className="text-xs text-gray-600 mt-1">{principle.article}</p>
                  <div className="mt-2 text-xs text-gray-500">
                    <div>Status: {principle.status}</div>
                    <div>Owner: {principle.owner}</div>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Data Processing Activities */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Data Processing Activities</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {processes.filter((p) => p.status === "active").length}
                </div>
                <div className="text-blue-600 text-sm">Active</div>
              </div>
              <div className="bg-gray-50 border border-gray-200 rounded p-4">
                <div className="text-gray-800 text-2xl font-bold">
                  {processes.filter((p) => p.status === "under_review").length}
                </div>
                <div className="text-gray-600 text-sm">Under Review</div>
              </div>
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {processes.filter((p) => p.status === "suspended").length}
                </div>
                <div className="text-red-600 text-sm">Suspended</div>
              </div>
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {processes.filter((p) => p.legal_basis === "consent").length}
                </div>
                <div className="text-green-600 text-sm">Consent-based</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Process</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Purpose</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Legal Basis</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Owner</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {processes.slice(0, 10).map((process) => (
                    <tr key={process.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {process.name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {process.purpose}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {process.legal_basis}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(process.status)}`}>
                          {process.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {process.owner}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Data Subject Requests */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Data Subject Requests</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {requests.filter((r) => r.status === "received").length}
                </div>
                <div className="text-blue-600 text-sm">Received</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {requests.filter((r) => r.status === "in_progress").length}
                </div>
                <div className="text-yellow-600 text-sm">In Progress</div>
              </div>
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {requests.filter((r) => r.status === "completed").length}
                </div>
                <div className="text-green-600 text-sm">Completed</div>
              </div>
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {requests.filter((r) => r.priority === "critical").length}
                </div>
                <div className="text-red-600 text-sm">Critical Priority</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Data Subject</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Priority</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Assigned To</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Due Date</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {requests.slice(0, 10).map((request) => (
                    <tr key={request.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {request.type}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {request.data_subject}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(request.status)}`}>
                          {request.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getPriorityColor(request.priority)}`}>
                          {request.priority}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {request.assigned_to}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(request.due_date).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Data Breaches */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Data Breaches</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {breaches.filter((b) => b.status === "open").length}
                </div>
                <div className="text-red-600 text-sm">Open</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {breaches.filter((b) => b.status === "investigating").length}
                </div>
                <div className="text-yellow-600 text-sm">Investigating</div>
              </div>
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {breaches.filter((b) => b.notification_required).length}
                </div>
                <div className="text-blue-600 text-sm">Notification Required</div>
              </div>
              <div className="bg-orange-50 border border-orange-200 rounded p-4">
                <div className="text-orange-800 text-2xl font-bold">
                  {breaches.reduce((sum, b) => sum + b.affected, 0)}
                </div>
                <div className="text-orange-600 text-sm">Total Affected</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Breach</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Affected</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Notification</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {breaches.slice(0, 10).map((breach) => (
                    <tr key={breach.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {breach.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {breach.type}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(breach.date).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {breach.affected}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(breach.status)}`}>
                          {breach.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {breach.notification_required ? (
                          <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800">
                            Required
                          </span>
                        ) : (
                          <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                            Not Required
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* DPIA */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Data Protection Impact Assessments</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {dpia.filter((d) => d.status === "required").length}
                </div>
                <div className="text-red-600 text-sm">Required</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {dpia.filter((d) => d.status === "in_progress").length}
                </div>
                <div className="text-yellow-600 text-sm">In Progress</div>
              </div>
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {dpia.filter((d) => d.status === "completed").length}
                </div>
                <div className="text-green-600 text-sm">Completed</div>
              </div>
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {dpia.filter((d) => d.status === "approved").length}
                </div>
                <div className="text-blue-600 text-sm">Approved</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">DPIA</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Controller</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Assessor</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Review Date</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {dpia.slice(0, 10).map((dpiaItem) => (
                    <tr key={dpiaItem.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {dpiaItem.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {dpiaItem.controller}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {dpiaItem.assessor}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(dpiaItem.status)}`}>
                          {dpiaItem.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(dpiaItem.review_date).toLocaleDateString()}
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
