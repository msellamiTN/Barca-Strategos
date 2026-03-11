package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barca-strategos/phoenix/pkg/types"
)

// ServiceNow connector creates and updates ServiceNow tickets.
type ServiceNow struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewServiceNow creates a ServiceNow connector.
func NewServiceNow(baseURL, username, password string) *ServiceNow {
	return &ServiceNow{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateTicket creates a ServiceNow incident.
func (sn *ServiceNow) CreateTicket(title, description, severity string) (string, error) {
	payload := map[string]interface{}{
		"short_description": title,
		"description":       description,
		"urgency":           mapSeverityToUrgency(severity),
		"state":             "1", // New
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", sn.BaseURL+"/api/now/table/incident", bytes.NewReader(body))
	req.SetBasicAuth(sn.Username, sn.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := sn.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("service now: %s", resp.Status)
	}
	var result struct {
		Result struct {
			SysID string `json:"sys_id"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Result.SysID, nil
}

// mapSeverityToUrgency converts severity to ServiceNow urgency.
func mapSeverityToUrgency(severity string) string {
	switch severity {
	case "critical":
		return "1"
	case "high":
		return "2"
	case "medium":
		return "3"
	case "low":
		return "4"
	default:
		return "3"
	}
}
