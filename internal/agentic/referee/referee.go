package referee

import (
	"context"
	"log"
)

// Referee reviews agent actions for safety/compliance.
type Referee struct{}

// New creates a Referee.
func New() *Referee {
	return &Referee{}
}

// Review evaluates a prompt/response pair.
func (r *Referee) Review(ctx context.Context, prompt, response string) (bool, string) {
	log.Println("referee: reviewing prompt/response")
	// TODO: policy checks, PII detection, compliance rules
	return true, "ok"
}
