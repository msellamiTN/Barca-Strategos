"use client";

import { useEffect, useState } from "react";

interface PCIRequirement {
  id: string;
  title: string;
  description: string;
  category: string;
  subcategories: string[];
  objective: string;
  control_type: string;
  status: string;
  implementation_date?: string;
  last_review_date?: string;
  evidence: string[];
  owner: string;
  risk_level: string;
}

interface PCIAssessment {
  assessment_id: string;
  timestamp: string;
  framework: string;
  version: string;
  scope: PCIScope;
  overall_score: number;
  requirement_assessments: PCIRequirementAssessment[];
  findings: PCIFinding[];
  recommendations: PCIRecommendation[];
  next_assessment_date: string;
}

interface PCIRequirementAssessment {
  requirement_id: string;
  requirement_title: string;
  category: string;
  compliance_score: number;
  status: string;
  findings: PCIFinding[];
  recommendations: string[];
  last_assessed: string;
}

interface PCIFinding {
  severity: string;
  description: string;
  recommendation: string;
  evidence_gaps: string[];
}

interface PCIRecommendation {
  priority: string;
  title: string;
  description: string;
  findings: PCIFinding[];
  estimated_effort: string;
  owner: string;
}

interface PCIScope {
  departments: string[];
  systems: string[];
  processes: string[];
}

interface PCIStats {
  total_requirements: number;
  compliant_requirements: number;
  implemented_requirements: number;
  partially_implemented_requirements: number;
  not_implemented_requirements: number;
  average_compliance_score: number;
  total_incidents: number;
  average_detection_time_minutes: number;
  average_response_time_minutes: number;
  last_incident_date?: string;
}

export default function PCIDashboard() {
  const [requirements, setRequirements] = useState<PCIRequirement[]>([]);
  const [selectedRequirement, setSelectedRequirement] = useState<PCIRequirement | null>(null);
  const [assessments, setAssessments] = useState<PCIAssessment[]>([]);
  const [latestAssessment, setLatestAssessment] = useState<PCIAssessment | null>(null);
  const [stats, setStats] = useState<PCIStats | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch requirements list
    fetch("/api/v1/compliance/pci_dss/requirements")
      .then((res) => res.json())
      .then((data) => setRequirements(data));

    // Fetch PCI DSS stats
    fetch("/api/v1/compliance/pci_dss/stats")
      .then((res) => res.json())
      .then((data) => setStats(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "pci_update") {
        setRequirements(msg.requirements);
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedRequirement) {
      // Fetch assessments
      fetch(`/api/v1/compliance/pci_dss/assessments`)
        .then((res) => res.json())
        .then((data) => {
          setAssessments(data);
          if (data.length > 0) {
            setLatestAssessment(data[data.length - 1]);
          }
        });
    }
  }, [selectedRequirement]);

  const conductAssessment = async () => {
    await fetch("/api/v1/compliance/pci_dss/assess", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        departments: ["Payment Processing", "Security"],
        systems: ["Payment Gateway", "Database"],
        processes: ["Card Processing"],
      }),
    });
    // Refresh assessments
    const res = await fetch("/api/v1/compliance/pci_dss/assessments");
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
      case "network_security": return "bg-red-50 border-red-200";
      case "system_configuration": return "bg-orange-50 border-orange-200";
      case "data_protection": return "bg-blue-50 border-blue-200";
      case "malware_protection": return "bg-green-50 border-green-200";
      case "secure_development": return "bg-purple-50 border-purple-200";
      case "access_control": return "bg-yellow-50 border-yellow-200";
      case "physical_security": return "bg-indigo-50 border-indigo-200";
      case "monitoring": return "bg-pink-50 border-pink-200";
      case "testing": return "bg-teal-50 border-teal-200";
      case "policy_management": return "bg-gray-50 border-gray-200";
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
      <h1 className="text-3xl font-bold mb-6">PCI DSS Compliance</h1>

      {/* PCI DSS Statistics */}
      {stats && (
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Compliance Overview</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 border border-blue-200 rounded p-4">
              <div className="text-blue-800 text-2xl font-bold">{stats.total_requirements}</div>
              <div className="text-blue-600 text-sm">Total Requirements</div>
            </div>
            <div className="bg-green-50 border border-green-200 rounded p-4">
              <div className="text-green-800 text-2xl font-bold">{stats.compliant_requirements}</div>
              <div className="text-green-600 text-sm">Compliant</div>
            </div>
            <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
              <div className="text-yellow-800 text-2xl font-bold">{stats.partially_implemented_requirements}</div>
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

      {/* Requirement Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={conductAssessment} className="px-4 py-2 bg-blue-600 text-white rounded">
            Conduct Assessment
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {requirements.map((requirement) => (
            <div
              key={requirement.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedRequirement?.id === requirement.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedRequirement(requirement)}
            >
              <h3 className="font-semibold">{requirement.title}</h3>
              <p className="text-sm text-gray-600 mb-2">{requirement.description}</p>
              <div className={`text-xs px-2 py-1 rounded mb-2 ${getCategoryColor(requirement.category)}`}>
                {requirement.category.replace("_", " ")}
              </div>
              <div className="flex justify-between items-center mb-2">
                <span className={`text-sm font-medium ${getStatusColor(requirement.status)}`}>
                  {requirement.status.replace("_", " ")}
                </span>
                <span className="text-xs text-gray-500">{requirement.risk_level}</span>
              </div>
              <div className="text-xs text-gray-500">
                Owner: {requirement.owner}
              </div>
            </div>
          ))}
        </div>
      </section>

      {selectedRequirement && (
        <>
          {/* Requirement Details */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Requirement Details</h2>
            <div className="border rounded-lg p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-semibold text-lg mb-2">{selectedRequirement.title}</h3>
                  <p className="text-gray-600 mb-4">{selectedRequirement.description}</p>
                  <div className="space-y-2">
                    <div><strong>Category:</strong> {selectedRequirement.category.replace("_", " ")}</div>
                    <div><strong>Type:</strong> {selectedRequirement.control_type.replace("_", " ")}</div>
                    <div><strong>Status:</strong> <span className={getStatusColor(selectedRequirement.status)}>{selectedRequirement.status.replace("_", " ")}</span></div>
                    <div><strong>Owner:</strong> {selectedRequirement.owner}</div>
                    <div><strong>Risk Level:</strong> {selectedRequirement.risk_level}</div>
                    {selectedRequirement.implementation_date && (
                      <div><strong>Implemented:</strong> {new Date(selectedRequirement.implementation_date).toLocaleDateString()}</div>
                    )}
                    {selectedRequirement.last_review_date && (
                      <div><strong>Last Review:</strong> {new Date(selectedRequirement.last_review_date).toLocaleDateString()}</div>
                    )}
                  </div>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Objective</h4>
                  <p className="text-sm text-gray-600 mb-4">{selectedRequirement.objective}</p>
                  <h4 className="font-semibold mb-2">Subcategories</h4>
                  <div className="space-y-1 mb-4">
                    {selectedRequirement.subcategories.map((subcategory, index) => (
                      <div key={index} className="text-sm text-gray-600">• {subcategory}</div>
                    ))}
                  </div>
                  <h4 className="font-semibold mb-2">Evidence</h4>
                  <div className="space-y-1">
                    {selectedRequirement.evidence.map((evidence, index) => (
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
                      Requirements Assessed: {latestAssessment.requirement_assessments.length}
                    </div>
                    <div className="text-sm text-gray-600">
                      Findings: {latestAssessment.findings.length}
                    </div>
                    <div className="text-sm text-gray-600">
                      Next Assessment: {new Date(latestAssessment.next_assessment_date).toLocaleDateString()}
                    </div>
                  </div>
                </div>

                {/* Requirement Assessment Details */}
                {latestAssessment.requirement_assessments.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Requirement Assessment Results</h4>
                    <div className="space-y-2">
                      {latestAssessment.requirement_assessments.map((assessment, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className="font-medium">{assessment.requirement_title}</span>
                            <span className={`ml-2 px-2 py-1 text-xs rounded ${getStatusColor(assessment.status)}`}>
                              {assessment.status.replace("_", " ")}
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
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Requirements</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Findings</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Recommendations</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Next Assessment</th>
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
                          {assessment.requirement_assessments.length}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.findings.length}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.recommendations.length}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {new Date(assessment.next_assessment_date).toLocaleDateString()}
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
