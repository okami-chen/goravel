package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	UserID     int    `json:"user_id"`
	TemplateID string `json:"template_id"`
}

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload, asynq.MaxRetry(100000)), nil
}

func NewEmailDeliveryTaskListner(ctx context.Context, t *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf(" [*] Successfully Task: %+v", payload)
	return nil
	//log.Printf(" [*] Attempting to Send Welcome Email to User %d...", p.UserID)
	//return fmt.Errorf("could not send email to the user") // <-- Return error
}
