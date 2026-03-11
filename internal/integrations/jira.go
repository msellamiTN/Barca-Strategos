package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barca-strategos/phoenix/pkg/types"
)

// Jira connector creates and updates Jira tickets.
type Jira struct {
	BaseURL    string
	Username   string
	Token      string
	HTTPClient *http.Client
}

// NewJira creates a Jira connector.
func NewJira(baseURL, username, token string) *Jira {
	return &Jira{
		BaseURL:  baseURL,
		Username: username,
		Token:    token,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateIssue creates a Jira issue.
func (j *Jira) CreateIssue(project, summary, description, priority string) (string, error) {
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{"key": project},
			"summary": summary,
			"description": map[string]interface{}{
				"type": "doc",
				"version": 1,
				"content": []map[string]interface{}{
					{"type": "paragraph", "content": []map[string]interface{}{{"type": "text", "text": description}}},
				},
			},
			"issuetype": map[string]string{"name": "Task"},
			"priority":  map[string]string{"name": priority},
		},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", j.BaseURL+"/rest/api/2/issue", bytes.NewReader(body))
	req.SetBasicAuth(j.Username, j.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := j.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("jira: %s", resp.Status)
	}
	var result struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Key, nil
}
