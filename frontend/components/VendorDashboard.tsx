"use client";

import { useEffect, useState } from "react";

interface Vendor {
  id: string;
  name: string;
  description: string;
  category: VendorCategory;
  contact_email: string;
  contact_phone?: string;
  address: string;
  services_offered: string[];
  contract_start_date: string;
  contract_end_date?: string;
  status: string;
  risk_level?: string;
  created_at: string;
  updated_at: string;
  last_assessed?: string;
}

interface VendorCategory {
  id: string;
  name: string;
  description: string;
  risk_factors: string[];
  base_risk_score: number;
  assessment_frequency_days: number;
}

interface VendorAssessment {
  assessment_id: string;
  vendor_id: string;
  vendor_name: string;
  timestamp: string;
  risk_assessment: VendorRiskAssessment;
  compliance_assessment: VendorComplianceAssessment;
  overall_risk_score: number;
  risk_level: string;
  recommendations: VendorRecommendation[];
  next_assessment_date: string;
}

interface VendorRiskAssessment {
  risk_score: number;
  risk_factors: RiskFactor[];
  financial_stability: number;
  operational_capability: number;
  security_posture: number;
  reputation_score: number;
}

interface RiskFactor {
  id: string;
  name: string;
  description: string;
  severity: string;
  score: number;
  mitigation_required: boolean;
}

interface VendorComplianceAssessment {
  compliance_score: number;
  compliance_issues: ComplianceIssue[];
  certifications: Certification[];
  regulatory_adherence: number;
  policy_compliance: number;
}

interface ComplianceIssue {
  id: string;
  description: string;
  severity: string;
  affected_regulations: string[];
  remediation_plan?: string;
}

interface Certification {
  name: string;
  issuer: string;
  obtained_date: string;
  expiry_date: string;
  status: string;
}

interface VendorRecommendation {
  priority: string;
  title: string;
  description: string;
  action_items: string[];
  owner: string;
  timeline: string;
}

interface VendorStats {
  total_vendors: number;
  active_vendors: number;
  critical_risk_vendors: number;
  high_risk_vendors: number;
  medium_risk_vendors: number;
  low_risk_vendors: number;
  average_risk_score: number;
  vendors_due_for_assessment: number;
  overdue_assessments: number;
}

export default function VendorDashboard() {
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [selectedVendor, setSelectedVendor] = useState<Vendor | null>(null);
  const [assessments, setAssessments] = useState<VendorAssessment[]>([]);
  const [latestAssessment, setLatestAssessment] = useState<VendorAssessment | null>(null);
  const [stats, setStats] = useState<VendorStats | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch vendors list
    fetch("/api/v1/compliance/vendor")
      .then((res) => res.json())
      .then((data) => setVendors(data));

    // Fetch vendor stats
    fetch("/api/v1/compliance/vendor/stats")
      .then((res) => res.json())
      .then((data) => setStats(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "vendor_update") {
        setVendors(msg.vendors);
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedVendor) {
      // Fetch assessments
      fetch(`/api/v1/compliance/vendor/${selectedVendor.id}/assessments`)
        .then((res) => res.json())
        .then((data) => {
          setAssessments(data);
          if (data.length > 0) {
            setLatestAssessment(data[data.length - 1]);
          }
        });
    }
  }, [selectedVendor]);

  const createVendor = async () => {
    await fetch("/api/v1/compliance/vendor", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: "Cloud Services Provider",
        description: "Enterprise cloud infrastructure and services provider",
        category_id: "CLOUD",
        contact_email: "contact@cloudprovider.com",
        address: "123 Tech Street, San Francisco, CA 94105",
        services_offered: ["Cloud Infrastructure", "Storage", "Computing", "Networking"],
        contract_start_date: new Date().toISOString(),
      }),
    });
    // Refresh list
    const res = await fetch("/api/v1/compliance/vendor");
    setVendors(await res.json());
  };

  const conductAssessment = async (vendorId: string) => {
    await fetch(`/api/v1/compliance/vendor/${vendorId}/assess`, {
      method: "POST",
    });
    // Refresh assessments
    if (selectedVendor && selectedVendor.id === vendorId) {
      const res = await fetch(`/api/v1/compliance/vendor/${vendorId}/assessments`);
      const data = await res.json();
      setAssessments(data);
      if (data.length > 0) {
        setLatestAssessment(data[data.length - 1]);
      }
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active": return "text-green-600";
      case "under_review": return "text-yellow-600";
      case "suspended": return "text-orange-600";
      case "terminated": return "text-red-600";
      default: return "text-gray-600";
    }
  };

  const getRiskLevelColor = (riskLevel: string) => {
    switch (riskLevel) {
      case "critical": return "bg-red-100 text-red-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "low": return "bg-green-100 text-green-800";
      case "minimal": return "bg-blue-100 text-blue-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const getCategoryColor = (categoryId: string) => {
    switch (categoryId) {
      case "CLOUD": return "bg-blue-50 border-blue-200";
      case "SOFTWARE": return "bg-green-50 border-green-200";
      case "CONSULTING": return "bg-purple-50 border-purple-200";
      case "INFRASTRUCTURE": return "bg-orange-50 border-orange-200";
      case "FINANCIAL": return "bg-red-50 border-red-200";
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

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">Vendor Risk Management</h1>

      {/* Vendor Statistics */}
      {stats && (
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Vendor Overview</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 border border-blue-200 rounded p-4">
              <div className="text-blue-800 text-2xl font-bold">{stats.total_vendors}</div>
              <div className="text-blue-600 text-sm">Total Vendors</div>
            </div>
            <div className="bg-green-50 border border-green-200 rounded p-4">
              <div className="text-green-800 text-2xl font-bold">{stats.active_vendors}</div>
              <div className="text-green-600 text-sm">Active</div>
            </div>
            <div className="bg-red-50 border border-red-200 rounded p-4">
              <div className="text-red-800 text-2xl font-bold">{stats.critical_risk_vendors + stats.high_risk_vendors}</div>
              <div className="text-red-600 text-sm">High Risk</div>
            </div>
            <div className="bg-purple-50 border border-purple-200 rounded p-4">
              <div className="text-purple-800 text-2xl font-bold">
                {stats.average_risk_score.toFixed(1)}
              </div>
              <div className="text-purple-600 text-sm">Avg Risk Score</div>
            </div>
          </div>
        </section>
      )}

      {/* Vendor Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={createVendor} className="px-4 py-2 bg-blue-600 text-white rounded">
            Add Vendor
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {vendors.map((vendor) => (
            <div
              key={vendor.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedVendor?.id === vendor.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedVendor(vendor)}
            >
              <h3 className="font-semibold">{vendor.name}</h3>
              <p className="text-sm text-gray-600 mb-2">{vendor.description}</p>
              <div className={`text-xs px-2 py-1 rounded mb-2 ${getCategoryColor(vendor.category.id)}`}>
                {vendor.category.name}
              </div>
              <div className="flex justify-between items-center mb-2">
                <span className={`text-sm font-medium ${getStatusColor(vendor.status)}`}>
                  {vendor.status.replace("_", " ")}
                </span>
                {vendor.risk_level && (
                  <span className={`text-xs px-2 py-1 rounded ${getRiskLevelColor(vendor.risk_level)}`}>
                    {vendor.risk_level}
                  </span>
                )}
              </div>
              <div className="text-xs text-gray-500">
                Contact: {vendor.contact_email}
              </div>
            </div>
          ))}
        </div>
      </section>

      {selectedVendor && (
        <>
          {/* Vendor Details */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Vendor Details</h2>
            <div className="border rounded-lg p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-semibold text-lg mb-2">{selectedVendor.name}</h3>
                  <p className="text-gray-600 mb-4">{selectedVendor.description}</p>
                  <div className="space-y-2">
                    <div><strong>Category:</strong> {selectedVendor.category.name}</div>
                    <div><strong>Status:</strong> <span className={getStatusColor(selectedVendor.status)}>{selectedVendor.status.replace("_", " ")}</span></div>
                    <div><strong>Contact:</strong> {selectedVendor.contact_email}</div>
                    {selectedVendor.contact_phone && (
                      <div><strong>Phone:</strong> {selectedVendor.contact_phone}</div>
                    )}
                    <div><strong>Address:</strong> {selectedVendor.address}</div>
                    <div><strong>Contract Start:</strong> {new Date(selectedVendor.contract_start_date).toLocaleDateString()}</div>
                    {selectedVendor.contract_end_date && (
                      <div><strong>Contract End:</strong> {new Date(selectedVendor.contract_end_date).toLocaleDateString()}</div>
                    )}
                    <div><strong>Created:</strong> {new Date(selectedVendor.created_at).toLocaleDateString()}</div>
                    {selectedVendor.last_assessed && (
                      <div><strong>Last Assessed:</strong> {new Date(selectedVendor.last_assessed).toLocaleDateString()}</div>
                    )}
                  </div>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Services Offered</h4>
                  <div className="space-y-1 mb-4">
                    {selectedVendor.services_offered.map((service, index) => (
                      <div key={index} className="text-sm text-gray-600">• {service}</div>
                    ))}
                  </div>
                  <h4 className="font-semibold mb-2">Risk Factors</h4>
                  <div className="space-y-1">
                    {selectedVendor.category.risk_factors.map((factor, index) => (
                      <div key={index} className="text-sm text-gray-600">• {factor}</div>
                    ))}
                  </div>
                  <div className="mt-4">
                    <div><strong>Base Risk Score:</strong> {selectedVendor.category.base_risk_score}</div>
                    <div><strong>Assessment Frequency:</strong> {selectedVendor.category.assessment_frequency_days} days</div>
                  </div>
                </div>
              </div>
              <div className="mt-4">
                <button
                  onClick={() => conductAssessment(selectedVendor.id)}
                  className="px-4 py-2 bg-green-600 text-white rounded"
                >
                  Conduct Assessment
                </button>
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
                    <div className="text-2xl font-bold text-orange-600">
                      {latestAssessment.overall_risk_score.toFixed(2)}
                    </div>
                    <div className="text-gray-600">Overall Risk Score</div>
                    <div className={`mt-2 inline-flex px-2 py-1 text-xs rounded ${getRiskLevelColor(latestAssessment.risk_level)}`}>
                      {latestAssessment.risk_level}
                    </div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-green-600">
                      {latestAssessment.compliance_assessment.compliance_score.toFixed(1)}%
                    </div>
                    <div className="text-gray-600">Compliance Score</div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">
                      Assessment Date: {new Date(latestAssessment.timestamp).toLocaleDateString()}
                    </div>
                    <div className="text-sm text-gray-600">
                      Next Assessment: {new Date(latestAssessment.next_assessment_date).toLocaleDateString()}
                    </div>
                  </div>
                </div>

                {/* Risk Assessment Details */}
                <div className="mb-6">
                  <h4 className="font-semibold mb-2">Risk Assessment Breakdown</h4>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="text-center">
                      <div className="text-lg font-semibold">{latestAssessment.risk_assessment.financial_stability.toFixed(2)}</div>
                      <div className="text-xs text-gray-600">Financial Stability</div>
                    </div>
                    <div className="text-center">
                      <div className="text-lg font-semibold">{latestAssessment.risk_assessment.operational_capability.toFixed(2)}</div>
                      <div className="text-xs text-gray-600">Operational Capability</div>
                    </div>
                    <div className="text-center">
                      <div className="text-lg font-semibold">{latestAssessment.risk_assessment.security_posture.toFixed(2)}</div>
                      <div className="text-xs text-gray-600">Security Posture</div>
                    </div>
                    <div className="text-center">
                      <div className="text-lg font-semibold">{latestAssessment.risk_assessment.reputation_score.toFixed(2)}</div>
                      <div className="text-xs text-gray-600">Reputation Score</div>
                    </div>
                  </div>
                </div>

                {/* Risk Factors */}
                {latestAssessment.risk_assessment.risk_factors.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Risk Factors</h4>
                    <div className="space-y-2">
                      {latestAssessment.risk_assessment.risk_factors.map((factor, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className={`px-2 py-1 text-xs rounded mr-2 ${getSeverityColor(factor.severity)}`}>
                              {factor.severity}
                            </span>
                            <span className="font-medium">{factor.name}</span>
                            <span className="text-sm text-gray-600 ml-2">({factor.description})</span>
                          </div>
                          <div className="text-sm">
                            Score: {factor.score.toFixed(2)}
                            {factor.mitigation_required && (
                              <span className="ml-2 text-red-600 text-xs">Mitigation Required</span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Compliance Issues */}
                {latestAssessment.compliance_assessment.compliance_issues.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Compliance Issues</h4>
                    <div className="space-y-2">
                      {latestAssessment.compliance_assessment.compliance_issues.map((issue, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className={`px-2 py-1 text-xs rounded mr-2 ${getSeverityColor(issue.severity)}`}>
                              {issue.severity}
                            </span>
                            {issue.description}
                          </div>
                          {issue.remediation_plan && (
                            <span className="text-xs text-blue-600">Has Remediation Plan</span>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Certifications */}
                {latestAssessment.compliance_assessment.certifications.length > 0 && (
                  <div className="mb-6">
                    <h4 className="font-semibold mb-2">Certifications</h4>
                    <div className="space-y-2">
                      {latestAssessment.compliance_assessment.certifications.map((cert, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <div className="font-medium">{cert.name}</div>
                            <div className="text-sm text-gray-600">Issued by: {cert.issuer}</div>
                          </div>
                          <div className="text-right">
                            <div className={`px-2 py-1 text-xs rounded ${
                              cert.status === "active" ? "bg-green-100 text-green-800" :
                              cert.status === "expired" ? "bg-red-100 text-red-800" :
                              "bg-gray-100 text-gray-800"
                            }`}>
                              {cert.status}
                            </div>
                            <div className="text-xs text-gray-600">
                              Expires: {new Date(cert.expiry_date).toLocaleDateString()}
                            </div>
                          </div>
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
                            <span className={`px-2 py-1 text-xs rounded ${
                              rec.priority === "critical" ? "bg-red-100 text-red-800" :
                              rec.priority === "high" ? "bg-orange-100 text-orange-800" :
                              rec.priority === "medium" ? "bg-yellow-100 text-yellow-800" :
                              "bg-gray-100 text-gray-800"
                            }`}>
                              {rec.priority}
                            </span>
                          </div>
                          <p className="text-sm text-gray-600 mb-2">{rec.description}</p>
                          <div className="text-sm text-gray-600">
                            <div><strong>Owner:</strong> {rec.owner}</div>
                            <div><strong>Timeline:</strong> {rec.timeline}</div>
                            <div><strong>Action Items:</strong></div>
                            <ul className="ml-4 list-disc">
                              {rec.action_items.map((item, itemIndex) => (
                                <li key={itemIndex}>{item}</li>
                              ))}
                            </ul>
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
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Risk Score</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Risk Level</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Compliance</th>
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
                          {assessment.overall_risk_score.toFixed(2)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRiskLevelColor(assessment.risk_level)}`}>
                            {assessment.risk_level}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {assessment.compliance_assessment.compliance_score.toFixed(1)}%
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
