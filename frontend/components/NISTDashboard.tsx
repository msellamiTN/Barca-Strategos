"use client";

import { useEffect, useState } from "react";

interface CSF {
  id: string;
  name: string;
  description: string;
  version: string;
  status: string;
  maturity_score?: number;
  created_at: string;
}

interface Function {
  id: string;
  name: string;
  description: string;
  category: string;
  status: string;
  progress: number;
  owner: string;
  next_review: string;
}

interface SubCategory {
  id: string;
  number: string;
  title: string;
  implementation_status: string;
  priority: string;
  owner: string;
  due_date: string;
}

interface Assessment {
  id: string;
  title: string;
  type: string;
  method: string;
  assessor: string;
  start_date: string;
  end_date: string;
  status: string;
  score: number;
}

export default function NISTDashboard() {
  const [csf, setCSF] = useState<CSF[]>([]);
  const [selectedCSF, setSelectedCSF] = useState<CSF | null>(null);
  const [functions, setFunctions] = useState<Function[]>([]);
  const [subcategories, setSubcategories] = useState<SubCategory[]>([]);
  const [assessments, setAssessments] = useState<Assessment[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch CSF list
    fetch("/api/v1/compliance/nist/csf")
      .then((res) => res.json())
      .then((data) => setCSF(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "nist_update") {
        setCSF(Object.values(msg.csf));
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedCSF) {
      // Fetch maturity score
      fetch(`/api/v1/compliance/nist/csf/${selectedCSF.id}/maturity`)
        .then((res) => res.json())
        .then((data) => {
          setSelectedCSF((prev) => prev ? { ...prev, maturity_score: data.maturity_score } : null);
        });

      // Fetch functions
      fetch("/api/v1/compliance/nist/functions")
        .then((res) => res.json())
        .then((data) => setFunctions(data));

      // Fetch subcategories
      fetch("/api/v1/compliance/nist/subcategories")
        .then((res) => res.json())
        .then((data) => setSubcategories(data));

      // Fetch assessments
      fetch("/api/v1/compliance/nist/assessments")
        .then((res) => res.json())
        .then((data) => setAssessments(data));
    }
  }, [selectedCSF]);

  const createCSF = async () => {
    await fetch("/api/v1/compliance/nist/csf", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Enterprise NIST CSF",
        description: "NIST Cybersecurity Framework Implementation",
        version: "1.1",
      }),
    });
    // Refresh list
    const res = await fetch("/api/v1/compliance/nist/csf");
    setCSF(await res.json());
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active": return "text-green-600";
      case "draft": return "text-gray-600";
      case "under_review": return "text-yellow-600";
      case "approved": return "text-blue-600";
      case "archived": return "text-red-600";
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

  const getFunctionColor = (category: string) => {
    switch (category) {
      case "identify": return "bg-blue-50 border-blue-200";
      case "protect": return "bg-green-50 border-green-200";
      case "detect": return "bg-yellow-50 border-yellow-200";
      case "respond": return "bg-orange-50 border-orange-200";
      case "recover": return "bg-purple-50 border-purple-200";
      default: return "bg-gray-50 border-gray-200";
    }
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">NIST Cybersecurity Framework</h1>

      {/* CSF Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={createCSF} className="px-4 py-2 bg-blue-600 text-white rounded">
            Create CSF
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {csf.map((csfItem) => (
            <div
              key={csfItem.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedCSF?.id === csfItem.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedCSF(csfItem)}
            >
              <h3 className="font-semibold">{csfItem.name}</h3>
              <p className="text-sm text-gray-600">{csfItem.description}</p>
              <div className="flex justify-between items-center mt-2">
                <span className={`text-sm font-medium ${getStatusColor(csfItem.status)}`}>
                  {csfItem.status}
                </span>
                <span className="text-sm text-gray-500">{csfItem.version}</span>
              </div>
              {csfItem.maturity_score !== undefined && (
                <div className="mt-2">
                  <div className="flex justify-between text-sm">
                    <span>Maturity Score</span>
                    <span className="font-medium">{csfItem.maturity_score.toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-green-600 h-2 rounded-full"
                      style={{ width: `${csfItem.maturity_score}%` }}
                    ></div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      {selectedCSF && (
        <>
          {/* Functions Overview */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Core Functions</h2>
            <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
              {functions.map((func) => (
                <div key={func.id} className={`border rounded-lg p-4 ${getFunctionColor(func.category)}`}>
                  <h3 className="font-semibold capitalize">{func.name}</h3>
                  <p className="text-sm text-gray-600 mt-1">{func.description}</p>
                  <div className="mt-3">
                    <div className="flex justify-between text-sm mb-1">
                      <span>Progress</span>
                      <span className="font-medium">{func.progress.toFixed(1)}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-blue-600 h-2 rounded-full"
                        style={{ width: `${func.progress}%` }}
                      ></div>
                    </div>
                  </div>
                  <div className="mt-2 text-xs text-gray-500">
                    <div>Status: {func.status}</div>
                    <div>Owner: {func.owner}</div>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Subcategories */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Implementation Subcategories</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {subcategories.filter((s) => s.implementation_status === "implemented").length}
                </div>
                <div className="text-green-600 text-sm">Implemented</div>
              </div>
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {subcategories.filter((s) => s.implementation_status === "tested").length}
                </div>
                <div className="text-blue-600 text-sm">Tested</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {subcategories.filter((s) => s.implementation_status === "partially_implemented").length}
                </div>
                <div className="text-yellow-600 text-sm">Partial</div>
              </div>
              <div className="bg-red-50 border border-red-200 rounded p-4">
                <div className="text-red-800 text-2xl font-bold">
                  {subcategories.filter((s) => s.implementation_status === "not_implemented").length}
                </div>
                <div className="text-red-600 text-sm">Not Implemented</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Subcategory</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Priority</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Owner</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Due Date</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {subcategories.slice(0, 10).map((subcat) => (
                    <tr key={subcat.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {subcat.number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {subcat.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(subcat.implementation_status)}`}>
                          {subcat.implementation_status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getPriorityColor(subcat.priority)}`}>
                          {subcat.priority}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {subcat.owner}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {new Date(subcat.due_date).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          {/* Assessments */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Assessments</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 border border-blue-200 rounded p-4">
                <div className="text-blue-800 text-2xl font-bold">
                  {assessments.filter((a) => a.status === "completed").length}
                </div>
                <div className="text-blue-600 text-sm">Completed</div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
                <div className="text-yellow-800 text-2xl font-bold">
                  {assessments.filter((a) => a.status === "in_progress").length}
                </div>
                <div className="text-yellow-600 text-sm">In Progress</div>
              </div>
              <div className="bg-green-50 border border-green-200 rounded p-4">
                <div className="text-green-800 text-2xl font-bold">
                  {assessments.filter((a) => a.status === "approved").length}
                </div>
                <div className="text-green-600 text-sm">Approved</div>
              </div>
              <div className="bg-gray-50 border border-gray-200 rounded p-4">
                <div className="text-gray-800 text-2xl font-bold">
                  {assessments.length > 0 ? (assessments.reduce((sum, a) => sum + a.score, 0) / assessments.length).toFixed(1) : "0"}
                </div>
                <div className="text-gray-600 text-sm">Avg Score</div>
              </div>
            </div>
            <div className="border rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Assessment</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Method</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Assessor</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Score</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {assessments.slice(0, 10).map((assessment) => (
                    <tr key={assessment.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {assessment.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {assessment.type}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {assessment.method}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {assessment.assessor}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(assessment.status)}`}>
                          {assessment.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {assessment.score.toFixed(1)}%
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
