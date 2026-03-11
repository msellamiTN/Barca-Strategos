package playbook

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Node represents a step in a playbook.
type Node struct {
	ID       string                 `json:"id"`
	Type     NodeType               `json:"type"`
	Config   map[string]interface{} `json:"config"`
	Next     []string               `json:"next,omitempty"`
}

// NodeType defines node types.
type NodeType string

const (
	NodeTypeTrigger   NodeType = "trigger"
	NodeTypeCondition NodeType = "condition"
	NodeTypeAction    NodeType = "action"
	NodeTypeDelay     NodeType = "delay"
)

// Playbook defines a workflow.
type Playbook struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Nodes       []Node  `json:"nodes"`
	StartNode   string  `json:"start_node"`
}

// Execution tracks a playbook run.
type Execution struct {
	ID        string                 `json:"id"`
	PlaybookID string                 `json:"playbook_id"`
	Status    ExecutionStatus        `json:"status"`
	StartTime time.Time              `json:"start_time"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Context   map[string]interface{} `json:"context"`
	CurrentNode string               `json:"current_node"`
	Results   []ExecutionResult      `json:"results"`
}

// ExecutionStatus represents run status.
type ExecutionStatus string

const (
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// ExecutionResult records node output.
type ExecutionResult struct {
	NodeID    string      `json:"node_id"`
	Status    string      `json:"status"`
	Output    interface{} `json:"output,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Engine executes playbooks.
type Engine struct {
	playbooks map[string]*Playbook
}

// New creates a Playbook engine.
func New() *Engine {
	return &Engine{
		playbooks: make(map[string]*Playbook),
	}
}

// RegisterPlaybook adds a playbook.
func (e *Engine) RegisterPlaybook(pb *Playbook) {
	e.playbooks[pb.ID] = pb
	log.Printf("playbook: registered %s", pb.ID)
}

// Execute starts a playbook execution.
func (e *Engine) Execute(ctx context.Context, playbookID string, initialContext map[string]interface{}) (*Execution, error) {
	pb, ok := e.playbooks[playbookID]
	if !ok {
		return nil, fmt.Errorf("playbook not found")
	}
	exec := &Execution{
		ID:         uuid.New().String(),
		PlaybookID: playbookID,
		Status:     ExecutionStatusRunning,
		StartTime:  time.Now().UTC(),
		Context:    initialContext,
		CurrentNode: pb.StartNode,
	}
	go e.run(ctx, exec, pb)
	return exec, nil
}

// run executes the playbook node graph.
func (e *Engine) run(ctx context.Context, exec *Execution, pb *Playbook) {
	for {
		select {
		case <-ctx.Done():
			exec.Status = ExecutionStatusCancelled
			return
		default:
		}
		node, ok := e.findNode(pb, exec.CurrentNode)
		if !ok {
			exec.Status = ExecutionStatusFailed
			return
		}
		result, next, err := e.executeNode(ctx, node, exec.Context)
		if err != nil {
			log.Printf("playbook: node %s failed: %v", node.ID, err)
			exec.Status = ExecutionStatusFailed
			return
		}
		exec.Results = append(exec.Results, ExecutionResult{
			NodeID:    node.ID,
			Status:    "completed",
			Output:    result,
			Timestamp: time.Now().UTC(),
		})
		if len(next) == 0 {
			exec.Status = ExecutionStatusCompleted
			now := time.Now().UTC()
			exec.EndTime = &now
			return
		}
		exec.CurrentNode = next[0]
	}
}

// findNode locates a node by ID.
func (e *Engine) findNode(pb *Playbook, nodeID string) (*Node, bool) {
	for _, n := range pb.Nodes {
		if n.ID == nodeID {
			return &n, true
		}
	}
	return nil, false
}

// executeNode runs a single node.
func (e *Engine) executeNode(ctx context.Context, node *Node, context map[string]interface{}) (interface{}, []string, error) {
	switch node.Type {
	case NodeTypeTrigger:
		return context, node.Next, nil
	case NodeTypeCondition:
		condition, _ := node.Config["condition"].(string)
		if e.evaluateCondition(condition, context) {
			return "true", node.Next, nil
		}
		return "false", []string{}, nil
	case NodeTypeAction:
		action, _ := node.Config["action"].(string)
		return e.executeAction(action, context), node.Next, nil
	case NodeTypeDelay:
		delay, _ := node.Config["seconds"].(float64)
		time.Sleep(time.Duration(delay) * time.Second)
		return "delayed", node.Next, nil
	default:
		return nil, nil, fmt.Errorf("unknown node type: %s", node.Type)
	}
}

// evaluateCondition evaluates a simple condition.
func (e *Engine) evaluateCondition(condition string, context map[string]interface{}) bool {
	// Simplified: if context has key matching condition, return true
	val, ok := context[condition]
	return ok && val != nil
}

// executeAction runs an action.
func (e *Engine) executeAction(action string, context map[string]interface{}) interface{} {
	log.Printf("playbook: executing action %s", action)
	return map[string]interface{}{"action": action, "result": "ok"}
}
