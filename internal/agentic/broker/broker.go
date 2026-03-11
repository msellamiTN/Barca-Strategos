package broker

import (
	"context"
	"fmt"
	"log"
)

// Tool represents a callable capability for agents.
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Inputs      map[string]string `json:"inputs"`
	SafetyTier  string            `json:"safetyTier"`
}

// Broker validates and executes agent tool calls.
type Broker struct {
	tools map[string]Tool
}

// New creates a Broker.
func New() *Broker {
	return &Broker{
		tools: make(map[string]Tool),
	}
}

// RegisterTool adds a tool to the registry.
func (b *Broker) RegisterTool(tool Tool) {
	b.tools[tool.Name] = tool
	log.Printf("broker: registered tool %s", tool.Name)
}

// ExecuteTool validates and executes a tool call.
func (b *Broker) ExecuteTool(ctx context.Context, toolName string, inputs map[string]interface{}) (map[string]interface{}, error) {
	tool, ok := b.tools[toolName]
	if !ok {
		return nil, fmt.Errorf("tool %s not found", toolName)
	}
	log.Printf("broker: executing tool %s (tier=%s)", tool.Name, tool.SafetyTier)
	// TODO: enforce safety tier, rate limits, audit
	return map[string]interface{}{"result": "ok"}, nil
}
