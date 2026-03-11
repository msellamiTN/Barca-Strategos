// Phoenix API Compliance Management System - Frontend Application

const API_BASE_URL = 'http://localhost:8080';

// State Management
const state = {
    currentPage: 'dashboard',
    apiEndpoint: API_BASE_URL,
    soc2Data: null,
    pciData: null,
    isLoading: false
};

// Initialize Application
document.addEventListener('DOMContentLoaded', () => {
    initializeNavigation();
    initializeEventListeners();
    loadDashboardData();
    checkAPIHealth();
});

// Navigation
function initializeNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    
    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const page = item.dataset.page;
            navigateToPage(page);
        });
    });
}

function navigateToPage(page) {
    // Update active nav item
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    document.querySelector(`[data-page="${page}"]`).classList.add('active');
    
    // Update page content
    document.querySelectorAll('.page').forEach(p => {
        p.classList.remove('active');
    });
    document.getElementById(`${page}-page`).classList.add('active');
    
    // Update page title
    const titles = {
        'dashboard': 'Compliance Dashboard',
        'soc2': 'SOC 2 Type II Compliance',
        'pci-dss': 'PCI DSS v4.0 Compliance',
        'reports': 'Compliance Reports',
        'settings': 'System Settings'
    };
    document.getElementById('page-title').textContent = titles[page];
    
    // Load page-specific data
    state.currentPage = page;
    loadPageData(page);
}

function loadPageData(page) {
    switch(page) {
        case 'dashboard':
            loadDashboardData();
            break;
        case 'soc2':
            loadSOC2Data();
            break;
        case 'pci-dss':
            loadPCIDSSData();
            break;
        case 'reports':
            loadReportsData();
            break;
    }
}

// Event Listeners
function initializeEventListeners() {
    // Refresh button
    document.getElementById('refresh-btn').addEventListener('click', () => {
        loadPageData(state.currentPage);
        showToast('Data refreshed', 'success');
    });
    
    // Run assessments
    document.getElementById('run-soc2-assessment').addEventListener('click', runSOC2Assessment);
    document.getElementById('run-pci-assessment').addEventListener('click', runPCIAssessment);
    
    // Generate report
    document.getElementById('generate-report').addEventListener('click', generateReport);
    
    // Save settings
    document.getElementById('save-settings').addEventListener('click', saveSettings);
}

// API Functions
async function apiRequest(endpoint, options = {}) {
    try {
        state.isLoading = true;
        const response = await fetch(`${state.apiEndpoint}${endpoint}`, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            }
        });
        
        if (!response.ok) {
            throw new Error(`API Error: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API Request Error:', error);
        showToast(`Error: ${error.message}`, 'error');
        return null;
    } finally {
        state.isLoading = false;
    }
}

async function checkAPIHealth() {
    const health = await apiRequest('/api/system/health');
    if (health && health.status === 'healthy') {
        document.querySelector('.status-dot').style.background = 'var(--success-color)';
        document.querySelector('.status-indicator span:last-child').textContent = 'System Online';
    } else {
        document.querySelector('.status-dot').style.background = 'var(--danger-color)';
        document.querySelector('.status-indicator span:last-child').textContent = 'System Offline';
    }
}

// Dashboard Functions
async function loadDashboardData() {
    // Load SOC2 stats
    const soc2Stats = await apiRequest('/api/v1/compliance/soc2/stats');
    if (soc2Stats) {
        document.getElementById('soc2-score').textContent = 
            `${Math.round(soc2Stats.compliance_score || 0)}%`;
        state.soc2Data = soc2Stats;
    }
    
    // Load PCI DSS stats
    const pciStats = await apiRequest('/api/v1/compliance/pci-dss/stats');
    if (pciStats) {
        document.getElementById('pci-score').textContent = 
            `${Math.round(pciStats.compliance_score || 0)}%`;
        state.pciData = pciStats;
    }
    
    // Calculate total findings and controls
    const totalFindings = (soc2Stats?.total_findings || 0) + (pciStats?.total_findings || 0);
    const totalControls = (soc2Stats?.total_controls || 0) + (pciStats?.total_requirements || 0);
    
    document.getElementById('total-findings').textContent = totalFindings;
    document.getElementById('total-controls').textContent = totalControls;
    
    // Update activity feed
    updateActivityFeed();
}

function updateActivityFeed() {
    const activities = [
        {
            icon: 'fa-check',
            color: 'success',
            title: 'SOC 2 Assessment Completed',
            time: '2 hours ago'
        },
        {
            icon: 'fa-exclamation-triangle',
            color: 'warning',
            title: 'New Finding Detected',
            time: '5 hours ago'
        },
        {
            icon: 'fa-file-alt',
            color: 'info',
            title: 'Compliance Report Generated',
            time: '1 day ago'
        }
    ];
    
    const activityList = document.getElementById('activity-list');
    activityList.innerHTML = activities.map(activity => `
        <div class="activity-item">
            <i class="fas ${activity.icon} text-${activity.color}"></i>
            <div>
                <p><strong>${activity.title}</strong></p>
                <span class="text-muted">${activity.time}</span>
            </div>
        </div>
    `).join('');
}

// SOC2 Functions
async function loadSOC2Data() {
    const container = document.getElementById('soc2-controls');
    container.innerHTML = '<div class="loading">Loading SOC 2 controls...</div>';
    
    const controls = await apiRequest('/api/v1/compliance/soc2/controls');
    
    if (controls && controls.length > 0) {
        container.innerHTML = controls.map(control => createControlCard(control)).join('');
    } else {
        container.innerHTML = '<div class="loading">No SOC 2 controls available. Run an assessment to get started.</div>';
    }
}

async function runSOC2Assessment() {
    showToast('Running SOC 2 assessment...', 'info');
    const result = await apiRequest('/api/v1/compliance/soc2/assess', { method: 'POST' });
    
    if (result) {
        showToast('SOC 2 assessment completed successfully', 'success');
        loadSOC2Data();
        loadDashboardData();
    }
}

// PCI DSS Functions
async function loadPCIDSSData() {
    const container = document.getElementById('pci-controls');
    container.innerHTML = '<div class="loading">Loading PCI DSS requirements...</div>';
    
    const requirements = await apiRequest('/api/v1/compliance/pci-dss/requirements');
    
    if (requirements && requirements.length > 0) {
        container.innerHTML = requirements.map(req => createControlCard(req, true)).join('');
    } else {
        container.innerHTML = '<div class="loading">No PCI DSS requirements available. Run an assessment to get started.</div>';
    }
}

async function runPCIAssessment() {
    showToast('Running PCI DSS assessment...', 'info');
    const result = await apiRequest('/api/v1/compliance/pci-dss/assess', { method: 'POST' });
    
    if (result) {
        showToast('PCI DSS assessment completed successfully', 'success');
        loadPCIDSSData();
        loadDashboardData();
    }
}

// Control Card Creation
function createControlCard(item, isPCI = false) {
    const id = isPCI ? item.requirement_id : item.control_id;
    const name = isPCI ? item.requirement_name : item.control_name;
    const description = item.description || 'No description available';
    const status = item.status || 'not_tested';
    const category = item.category || 'General';
    
    const statusClass = {
        'compliant': 'compliant',
        'non_compliant': 'non-compliant',
        'partial': 'partial',
        'not_tested': 'partial'
    }[status] || 'partial';
    
    const statusText = {
        'compliant': 'Compliant',
        'non_compliant': 'Non-Compliant',
        'partial': 'Partial',
        'not_tested': 'Not Tested'
    }[status] || 'Unknown';
    
    return `
        <div class="control-card ${statusClass}">
            <div class="control-header">
                <div class="control-id">${id}</div>
                <span class="control-status ${statusClass}">${statusText}</span>
            </div>
            <h4>${name}</h4>
            <p class="control-description">${description}</p>
            <div class="control-footer">
                <span><i class="fas fa-tag"></i> ${category}</span>
                <span><i class="fas fa-calendar"></i> Last checked: Today</span>
            </div>
        </div>
    `;
}

// Reports Functions
async function loadReportsData() {
    const container = document.getElementById('reports-list');
    container.innerHTML = '<div class="loading">Loading reports...</div>';
    
    // Mock reports data
    const reports = [
        {
            id: 1,
            title: 'SOC 2 Type II Compliance Report',
            date: '2026-03-11',
            type: 'SOC2',
            status: 'completed'
        },
        {
            id: 2,
            title: 'PCI DSS v4.0 Assessment Report',
            date: '2026-03-10',
            type: 'PCI-DSS',
            status: 'completed'
        },
        {
            id: 3,
            title: 'Monthly Compliance Summary',
            date: '2026-03-01',
            type: 'Summary',
            status: 'completed'
        }
    ];
    
    container.innerHTML = `
        <div class="controls-grid">
            ${reports.map(report => `
                <div class="control-card">
                    <div class="control-header">
                        <div class="control-id">${report.type}</div>
                        <span class="control-status compliant">${report.status}</span>
                    </div>
                    <h4>${report.title}</h4>
                    <p class="control-description">Generated on ${report.date}</p>
                    <div class="control-footer">
                        <button class="btn btn-primary" onclick="downloadReport(${report.id})">
                            <i class="fas fa-download"></i> Download
                        </button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;
}

async function generateReport() {
    showToast('Generating compliance report...', 'info');
    
    const soc2Report = await apiRequest('/api/v1/compliance/soc2/report', { method: 'POST' });
    const pciReport = await apiRequest('/api/v1/compliance/pci-dss/report', { method: 'POST' });
    
    if (soc2Report || pciReport) {
        showToast('Reports generated successfully', 'success');
        loadReportsData();
    }
}

function downloadReport(reportId) {
    showToast(`Downloading report ${reportId}...`, 'info');
    // Implement actual download logic
}

// Settings Functions
function saveSettings() {
    const apiEndpoint = document.getElementById('api-endpoint').value;
    const monitoringInterval = document.getElementById('monitoring-interval').value;
    
    state.apiEndpoint = apiEndpoint;
    
    showToast('Settings saved successfully', 'success');
}

// Toast Notification
function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// Export functions for inline use
window.downloadReport = downloadReport;
