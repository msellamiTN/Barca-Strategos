package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/barca-strategos/phoenix/internal/agentic/broker"
)

// SlackBot handles Slack interactions.
type SlackBot struct {
	broker *broker.Broker
}

// NewSlackBot creates a Slack bot.
func NewSlackBot(b *broker.Broker) *SlackBot {
	return &SlackBot{broker: b}
}

// HandleSlashCommand processes Slack slash commands.
func (b *SlackBot) HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cmd := r.FormValue("command")
	text := r.FormValue("text")
	switch cmd {
	case "/phoenix-status":
		b.handleStatus(w, text)
	case "/phoenix-risk":
		b.handleRisk(w, text)
	case "/phoenix-alert":
		b.handleAlert(w, text)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unknown command"))
	}
}

func (b *SlackBot) handleStatus(w http.ResponseWriter, text string) {
	result, err := b.broker.ExecuteTool(context.Background(), "resolve_alert", map[string]interface{}{})
	if err != nil {
		log.Printf("slack: tool error: %v", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error executing tool"))
		return
	}
	resp := map[string]interface{}{"response_type": "in_channel", "text": "Status OK", "details": result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (b *SlackBot) handleRisk(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Risk feature coming soon"))
}

func (b *SlackBot) handleAlert(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert feature coming soon"))
}
