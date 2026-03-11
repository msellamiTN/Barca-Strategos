"use client";

import { useEffect, useState } from "react";

interface Policy {
  id: string;
  title: string;
  description: string;
  category: PolicyCategory;
  version: string;
  status: string;
  owner: string;
  approvers: string[];
  created_at: string;
  updated_at: string;
  published_at?: string;
  review_date?: string;
  tags: string[];
}

interface PolicyCategory {
  id: string;
  name: string;
  description: string;
  subcategories: string[];
  approval_required: boolean;
  review_period_days: number;
}

interface PolicyComplianceStatus {
  policy_id: string;
  policy_title: string;
  compliance_score: number;
  last_checked: string;
  compliance_issues: ComplianceIssue[];
  next_review_date: string;
}

interface ComplianceIssue {
  id: string;
  description: string;
  severity: string;
  affected_departments: string[];
  remediation_required: boolean;
}

interface ApprovalWorkflow {
  id: string;
  policy_id: string;
  approvers: Approver[];
  created_at: string;
  expires_at: string;
  status: string;
}

interface Approver {
  id: string;
  name: string;
  email: string;
  role: string;
  decision?: string;
  decision_date?: string;
  comments?: string;
}

interface PolicyStats {
  total_policies: number;
  draft_policies: number;
  pending_approval: number;
  published_policies: number;
  archived_policies: number;
  average_compliance_score: number;
  policies_due_for_review: number;
  overdue_reviews: number;
}

export default function PolicyDashboard() {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [selectedPolicy, setSelectedPolicy] = useState<Policy | null>(null);
  const [complianceStatus, setComplianceStatus] = useState<PolicyComplianceStatus | null>(null);
  const [workflows, setWorkflows] = useState<ApprovalWorkflow[]>([]);
  const [stats, setStats] = useState<PolicyStats | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Fetch policies list
    fetch("/api/v1/compliance/policy")
      .then((res) => res.json())
      .then((data) => setPolicies(data));

    // Fetch policy stats
    fetch("/api/v1/compliance/policy/stats")
      .then((res) => res.json())
      .then((data) => setStats(data));

    // WebSocket for real-time updates
    const socket = new WebSocket(`ws://${window.location.host}/ws`);
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "policy_update") {
        setPolicies(msg.policies);
      }
    };
    setWs(socket);
    return () => socket.close();
  }, []);

  useEffect(() => {
    if (selectedPolicy) {
      // Fetch compliance status
      fetch(`/api/v1/compliance/policy/${selectedPolicy.id}/compliance`)
        .then((res) => res.json())
        .then((data) => setComplianceStatus(data));

      // Fetch workflows
      fetch(`/api/v1/compliance/policy/${selectedPolicy.id}/workflows`)
        .then((res) => res.json())
        .then((data) => setWorkflows(data));
    }
  }, [selectedPolicy]);

  const createPolicy = async () => {
    await fetch("/api/v1/compliance/policy", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        title: "Enterprise Security Policy",
        description: "Comprehensive security policy for enterprise operations",
        content: "This policy establishes security requirements for all enterprise operations...",
        category_id: "SECURITY",
        owner: "security@company.com",
        approvers: ["policy_committee"],
        tags: ["security", "enterprise"],
        priority: "high",
      }),
    });
    // Refresh list
    const res = await fetch("/api/v1/compliance/policy");
    setPolicies(await res.json());
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "published": return "text-green-600";
      case "draft": return "text-gray-600";
      case "pending_approval": return "text-yellow-600";
      case "approved": return "text-blue-600";
      case "under_review": return "text-purple-600";
      case "archived": return "text-red-600";
      default: return "text-gray-600";
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "low": return "bg-green-100 text-green-800";
      case "medium": return "bg-yellow-100 text-yellow-800";
      case "high": return "bg-orange-100 text-orange-800";
      case "critical": return "bg-red-100 text-red-800";
      default: return "bg-gray-100 text-gray-800";
    }
  };

  const getCategoryColor = (categoryId: string) => {
    switch (categoryId) {
      case "SECURITY": return "bg-blue-50 border-blue-200";
      case "COMPLIANCE": return "bg-green-50 border-green-200";
      case "OPERATIONAL": return "bg-purple-50 border-purple-200";
      case "HR": return "bg-orange-50 border-orange-200";
      default: return "bg-gray-50 border-gray-200";
    }
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">Policy Management</h1>

      {/* Policy Statistics */}
      {stats && (
        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Policy Overview</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-blue-50 border border-blue-200 rounded p-4">
              <div className="text-blue-800 text-2xl font-bold">{stats.total_policies}</div>
              <div className="text-blue-600 text-sm">Total Policies</div>
            </div>
            <div className="bg-green-50 border border-green-200 rounded p-4">
              <div className="text-green-800 text-2xl font-bold">{stats.published_policies}</div>
              <div className="text-green-600 text-sm">Published</div>
            </div>
            <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
              <div className="text-yellow-800 text-2xl font-bold">{stats.pending_approval}</div>
              <div className="text-yellow-600 text-sm">Pending Approval</div>
            </div>
            <div className="bg-purple-50 border border-purple-200 rounded p-4">
              <div className="text-purple-800 text-2xl font-bold">
                {stats.average_compliance_score.toFixed(1)}%
              </div>
              <div className="text-purple-600 text-sm">Avg Compliance</div>
            </div>
          </div>
        </section>
      )}

      {/* Policy Selection */}
      <section className="mb-6">
        <div className="flex gap-4 mb-4">
          <button onClick={createPolicy} className="px-4 py-2 bg-blue-600 text-white rounded">
            Create Policy
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {policies.map((policy) => (
            <div
              key={policy.id}
              className={`border rounded-lg p-4 cursor-pointer hover:bg-gray-50 ${
                selectedPolicy?.id === policy.id ? "border-blue-500 bg-blue-50" : "border-gray-200"
              }`}
              onClick={() => setSelectedPolicy(policy)}
            >
              <h3 className="font-semibold">{policy.title}</h3>
              <p className="text-sm text-gray-600 mb-2">{policy.description}</p>
              <div className={`text-xs px-2 py-1 rounded mb-2 ${getCategoryColor(policy.category.id)}`}>
                {policy.category.name}
              </div>
              <div className="flex justify-between items-center">
                <span className={`text-sm font-medium ${getStatusColor(policy.status)}`}>
                  {policy.status.replace("_", " ")}
                </span>
                <span className="text-sm text-gray-500">v{policy.version}</span>
              </div>
              <div className="mt-2 text-xs text-gray-500">
                Owner: {policy.owner}
              </div>
            </div>
          ))}
        </div>
      </section>

      {selectedPolicy && (
        <>
          {/* Policy Details */}
          <section className="mb-8">
            <h2 className="text-xl font-semibold mb-4">Policy Details</h2>
            <div className="border rounded-lg p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h3 className="font-semibold text-lg mb-2">{selectedPolicy.title}</h3>
                  <p className="text-gray-600 mb-4">{selectedPolicy.description}</p>
                  <div className="space-y-2">
                    <div><strong>Category:</strong> {selectedPolicy.category.name}</div>
                    <div><strong>Version:</strong> {selectedPolicy.version}</div>
                    <div><strong>Owner:</strong> {selectedPolicy.owner}</div>
                    <div><strong>Status:</strong> <span className={getStatusColor(selectedPolicy.status)}>{selectedPolicy.status.replace("_", " ")}</span></div>
                    <div><strong>Created:</strong> {new Date(selectedPolicy.created_at).toLocaleDateString()}</div>
                    <div><strong>Updated:</strong> {new Date(selectedPolicy.updated_at).toLocaleDateString()}</div>
                    {selectedPolicy.published_at && (
                      <div><strong>Published:</strong> {new Date(selectedPolicy.published_at).toLocaleDateString()}</div>
                    )}
                    {selectedPolicy.review_date && (
                      <div><strong>Review Date:</strong> {new Date(selectedPolicy.review_date).toLocaleDateString()}</div>
                    )}
                  </div>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Approvers</h4>
                  <div className="space-y-1">
                    {selectedPolicy.approvers.map((approver, index) => (
                      <div key={index} className="text-sm text-gray-600">• {approver}</div>
                    ))}
                  </div>
                  {selectedPolicy.tags.length > 0 && (
                    <>
                      <h4 className="font-semibold mb-2 mt-4">Tags</h4>
                      <div className="flex flex-wrap gap-2">
                        {selectedPolicy.tags.map((tag, index) => (
                          <span key={index} className="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs">
                            {tag}
                          </span>
                        ))}
                      </div>
                    </>
                  )}
                </div>
              </div>
            </div>
          </section>

          {/* Compliance Status */}
          {complianceStatus && (
            <section className="mb-8">
              <h2 className="text-xl font-semibold mb-4">Compliance Status</h2>
              <div className="border rounded-lg p-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <div className="text-2xl font-bold text-green-600">
                      {complianceStatus.compliance_score.toFixed(1)}%
                    </div>
                    <div className="text-gray-600">Compliance Score</div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">
                      Last Checked: {new Date(complianceStatus.last_checked).toLocaleDateString()}
                    </div>
                    <div className="text-sm text-gray-600">
                      Next Review: {new Date(complianceStatus.next_review_date).toLocaleDateString()}
                    </div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600">
                      Issues: {complianceStatus.compliance_issues.length}
                    </div>
                  </div>
                </div>
                
                {complianceStatus.compliance_issues.length > 0 && (
                  <div className="mt-4">
                    <h4 className="font-semibold mb-2">Compliance Issues</h4>
                    <div className="space-y-2">
                      {complianceStatus.compliance_issues.map((issue, index) => (
                        <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                          <div>
                            <span className={`px-2 py-1 text-xs rounded mr-2 ${getSeverityColor(issue.severity)}`}>
                              {issue.severity}
                            </span>
                            {issue.description}
                          </div>
                          {issue.remediation_required && (
                            <span className="text-xs text-red-600">Remediation Required</span>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </section>
          )}

          {/* Approval Workflows */}
          {workflows.length > 0 && (
            <section className="mb-8">
              <h2 className="text-xl font-semibold mb-4">Approval Workflows</h2>
              <div className="border rounded-lg overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Workflow</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Approvers</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Expires</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {workflows.map((workflow) => (
                      <tr key={workflow.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {workflow.id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <div className="space-y-1">
                            {workflow.approvers.map((approver, index) => (
                              <div key={index} className="flex items-center space-x-2">
                                <span>{approver.name}</span>
                                {approver.decision && (
                                  <span className={`px-2 py-1 text-xs rounded ${
                                    approver.decision === "approved" ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
                                  }`}>
                                    {approver.decision}
                                  </span>
                                )}
                              </div>
                            ))}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                            workflow.status === "approved" ? "bg-green-100 text-green-800" :
                            workflow.status === "pending" ? "bg-yellow-100 text-yellow-800" :
                            workflow.status === "rejected" ? "bg-red-100 text-red-800" :
                            "bg-gray-100 text-gray-800"
                          }`}>
                            {workflow.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {new Date(workflow.created_at).toLocaleDateString()}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {new Date(workflow.expires_at).toLocaleDateString()}
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
