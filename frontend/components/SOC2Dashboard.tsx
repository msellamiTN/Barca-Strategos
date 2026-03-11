"use client";

import { useEffect, useState } from "react";

interface SOC2Control {
  id: string;
  title: string;
  description: string;
  category: SOC2ControlCategory;
  subcategories: string[];
  objective: string;
  control_type: SOC2ControlType;
  status: SOC2ControlStatus;
  implementation_date?: string;
  last_review_date?: string;
  evidence: string[];
  owner: string;
  risk_level: string;
}

interface SOC2ControlCategory {
  value: string;
  label: string;
}

interface SOC2ControlType {
  value: string;
  label: string;
}

interface SOC2ControlStatus {
  value: string;
  label: string;
}

interface SOC2Assessment {
  assessment_id: string;
  timestamp: string;
  framework: string;
  version: string;
  scope: SOC2Scope;
  overall_score: number;
  control_assessments: SOC2ControlAssessment[];
  findings: SOC2Finding[];
  recommendations: SOC2Recommendation[];
  last_assessed: string;
}

interface SOC2ControlAssessment {
  control_id: string;
  control_title: string;
  category: SOC2ControlCategory;
  compliance_score: number;
  status: SOC2ControlStatus;
  findings: SOC2Finding[];
  recommendations: string[];
  last_assessed: string;
}

interface SOC2Finding {
  severity: string;
  description: string;
  recommendation: string;
  evidence_gaps: string[];
}

interface SOC2Recommendation {
  priority: string;
  title: string;
  description: string;
  findings: SOC2Finding[];
  estimated_effort: string;
  owner: string;
}

interface SOC2Scope {
  departments: string[];
  systems: string[];
  processes: string[];
}

interface SOC2Stats {
  total_controls: number;
  compliant_controls: number;
  implemented_controls: number;
  partially_implemented_controls: number;
  not_implemented_controls: number;
  average_compliance_score: number;
  total_incidents: number;
  average_detection_time_minutes: number;
  average_response_time_minutes: number;
  last_incident_date?: string;
}

export default function SOC2Dashboard() {
  const [controls, setControls] = useState<SOC2Control[]>([]);
  const [selectedControl, setSelectedControl] = useState<SOC2Control | null>(null);
  const [assessments, setAssessments] = useState<SOC2Assessment[]>([]);
  const [latestAssessment, setLatestAssessment] = useState<SOC2Assessment | null>(null);
  const [stats, setStats] = useState<SOC2Stats | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch controls list
    fetch("/api/v1/compliance/soc2/controls")
      .then((res) => res.json())
      .then((data) => setControls(data));

    // Fetch SOC2 stats
    fetch("/api/v1/compliance/soc2/stats")
      .then((res) => res.json())
      .then((data) => setStats(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "soc2_update") {
        setControls(msg.controls);
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedControl) {
      // Fetch assessments
      fetch(`/api/v1/compliance/soc2/assessments`)
        .then((res) => res.json())
        .then((data) => {
          setAssessments(data);
          if (data.length > 0) {
            setLatestAssessment(data[data.length - 1]);
          }
        });
    }
  }, [selectedControl]);

  const conductAssessment = async () => {
    await fetch("/api/v1/compliance/soc2/assess", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        departments: ["IT", "Security"],
        systems: ["Phoenix Core", "Database"],
        processes: ["Incident Response"],
      }),
    });
    // Refresh assessments
    const res = await fetch("/api/v1/compliance/soc2/assessments");
    setAssessments(await res.json());
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "compliant": return "text-green-600";
      case "implemented": return "text-blue-600";
      case "partially_implemented": return "text-yellow-600";
      case "not_implemented": return "text-red-600";
      default: return "text-gray-600";
    }
  };

  const getCategoryColor = (category: string) => {
    switch (category) {
      case "governance": return "bg-purple-50 border-purple-200";
      case "asset_management": return "bg-blue-50 border-blue-200";
      case "access_control": return "bg-green-50 border-green-200";
      case "operational": return "bg-orange-50 border-orange-200";
      case "incident_response": return "bg-red-50 border-red-200";
      case "vulnerability_management": return "bg-yellow-50 border-yellow-200";
      case "disaster_recovery": return "bg-indigo-50 border-indigo-200";
      case "test_evaluation": return "bg-pink-50 border-pink-200";
      case "communications_security": return "bg-teal-50 border-teal-200";
      default: return "bg-gray-50 border-gray-200";
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "critical": return "bg-red-100 text-red-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "low": return "bg-green-100 text-green-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "critical": return "bg-red-100 text-red-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "low": return "bg-green-100 text-green-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">SOC 2 Type II Compliance</h1>

      {/* SOC2 Statistics */}
      {stats && (
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Compliance Overview</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 border border-blue-200 rounded p-4">
              <div className="text-blue-800 text-2xl font-bold">{stats.total_controls}</div>
              <div className="text-blue-600 text-sm">Total Controls</div>
            </div>
            <div className="bg-green-50 border border-green-200 rounded p-4">
              <div className="text-green-800 text-2xl font-bold">{stats.compliant_controls}</div>
              <div className="text-green-600 text-sm">Compliant</div>
            </div>
            <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
              <div className="text-yellow-800 text-2xl font-bold">{stats.partially_implemented_controls}</div>
              <div className="text-yellow-600 text-sm">Partial</div>
            </div>
            <div className="bg-purple-50 border border-purple-200 rounded p-4">
              <div className="text-purple-800 text-2xl font-bold">
                {(stats.average_compliance_score * 100).toFixed(1)}%
              </div>
              <div className="text-purple-600 text-sm">Avg Score</div>
            </div>
          </div>
        </section>
      )}

      {/* Control Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={conductAssessment} className="px-4 py-2 bg-blue-600 text-white rounded">
            Conduct Assessment
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {controls.map((control) => (
            <div
              key={control.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedControl?.id === control.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedControl(control)}
            >
              <h3 className="font-semibold">{control.title}</h3>
              <p className="text-sm text-gray-600 mb-2">{control.description}</p>
              <div className={`text-xs px-2 py-1 rounded mb-2 ${getCategoryColor(control.category)}`}>
                {control.category.replace("_", " ")}
              </div>
              <div className="flex justify-between items-center mb-2">
                <span className={`text-sm font-medium ${getStatusColor(control.status)}`}>
                  {control.status.replace("_", " ")}
                </span>
                <span className="text-xs text-gray-500">{control.risk_level}</span>
              </div>
              <div className="text-xs text-gray-500">
                Owner: {control.owner}
              </div>
            </div>
          ))}
        </div>
      </section>

      {selectedControl && (
        <>
          {/* Control Details */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Control Details</h2>
            <div className="border rounded-lg p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-semibold text-lg mb-2">{selectedControl.title}</h3>
                  <p className="text-gray-600 mb-4">{selectedControl.description}</p>
                  <div className="space-y-2">
                    <div><strong>Category:</strong> {selectedControl.category.replace("_", " ")}</div>
                    <div><strong>Type:</strong> {selectedControl.control_type.replace("_", " ")}</div>
                    <div><strong>Status:</strong> <span className={getStatusColor(selectedControl.status)}>{selectedControl.status.replace("_", " ")}</span></div>
                    <div><strong>Owner:</strong> {selectedControl.owner}</div>
                    <div><strong>Risk Level:</strong> {selectedControl.risk_level}</div>
                    {selectedControl.implementation_date && (
                      <div><strong>Implemented:</strong> {new Date(selectedControl.implementation_date).toLocaleDateString()}</div>
                    )}
                    {selectedControl.last_review_date && (
                      <div><strong>Last Review:</strong> {new Date(selectedControl.last_review_date).toLocaleDateString()}</div>
                    )}
                  </div>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Objective</h4>
                  <p className="text-sm text-gray-600 mb-4">{selectedControl.objective}</p>
                  <h4 className="font-semibold mb-2">Subcategories</h4>
                  <div className="space-y-1 mb-4">
                    {selectedControl.subcategories.map((subcategory, index) => (
                      <div key={index} className="text-sm text-gray-600">• {subcategory}</div>
                    ))}
                  </div>
                  <h4 className="font-semibold mb-2">Evidence</h4>
                  <div className="space-y-1">
                    {selectedControl.evidence.map((evidence, index) => (
                      <div key={index} className="text-sm text-gray-600">• {evidence}</div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </section>

          {/* Latest Assessment */}
          {latestAssessment && (
            <section className="mb-8">
              <h2 className="text-xl font-semibold mb-4">Latest Assessment</h2>
              <div className="border rounded-lg p-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                  <div>
                    <div className="text-2xl font-bold text-blue-600">
                      {(latestAssessment.overall_score * 100).toFixed(1)}%
                    </div>
                    <div className="text-gray-600">Overall Compliance Score</div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">
                      Assessment Date: {new Date(latestAssessment.timestamp).toLocaleDateString()}
                    </div>
                    <div className="text-sm text-gray-600">
                      Framework: {latestAssessment.framework} {latestAssessment.version}
                    </div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">
                      Controls Assessed: {latestAssessment.control_assessments.length}
                    </div>
                    <div className="text-sm text-gray-600">
                      Findings: {latestAssessment.findings.length}
                    </div>
                  </div>
                </div>

                {/* Control Assessment Details */}
                {latestAssessment.control_assessments.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Control Assessment Results</h4>
                    <div className="space-y-2">
                      {latestAssessment.control_assessments.map((assessment, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className="font-medium">{assessment.control_title}</span>
                            <span className={`ml-2 px-2 py-1 text-xs rounded ${getStatusColor(assessment.status.value)}`}>
                              {assessment.status.label}
                            </span>
                          </div>
                          <div className="text-sm">
                            Score: {(assessment.compliance_score * 100).toFixed(1)}%
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Findings */}
                {latestAssessment.findings.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Findings</h4>
                    <div className="space-y-2">
                      {latestAssessment.findings.map((finding, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className={`px-2 py-1 text-xs rounded mr-2 ${getSeverityColor(finding.severity)}`}>
                              {finding.severity}
                            </span>
                            {finding.description}
                          </div>
                          {finding.evidence_gaps.length > 0 && (
                            <span className="text-xs text-blue-600">
                              {finding.evidence_gaps.length} evidence gaps
                            </span>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Recommendations */}
                {latestAssessment.recommendations.length > 0 && (
                  <div>
                    <h4 className="font-semibold mb-2">Recommendations</h4>
                    <div className="space-y-3">
                      {latestAssessment.recommendations.map((rec, index) => (
                        <div key={index} className="border rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <h5 className="font-semibold">{rec.title}</h5>
                            <span className={`px-2 py-1 text-xs rounded ${getPriorityColor(rec.priority)}`}>
                              {rec.priority}
                            </span>
                          </div>
                          <p className="text-sm text-gray-600 mb-2">{rec.description}</p>
                          <div className="text-sm text-gray-600">
                            <div><strong>Owner:</strong> {rec.owner}</div>
                            <div><strong>Estimated Effort:</strong> {rec.estimated_effort}</div>
                            <div><strong>Findings:</strong> {rec.findings.length}</div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </section>
          )}

          {/* Assessment History */}
          {assessments.length > 1 && (
            <section className="mb-8">
              <h2 className="text-xl font-semibold mb-4">Assessment History</h2>
              <div className="border rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Score</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Controls</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Findings</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Recommendations</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {assessments.map((assessment) => (
                      <tr key={assessment.assessment_id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {new Date(assessment.timestamp).toLocaleDateString()}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {(assessment.overall_score * 100).toFixed(1)}%
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.control_assessments.length}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.findings.length}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.recommendations.length}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </section>
          )}
        </>
      )}
    </main>
  );
}
