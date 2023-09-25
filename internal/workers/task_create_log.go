package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "go-k8s/internal/db/sqlc"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskCreateLoginLog = "task:create_login_log"

type PayloadCreateLoginLog struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskCreateLoginLog(
	ctx context.Context,
	payload *PayloadCreateLoginLog,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskCreateLoginLog, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskCreateLoginLog(ctx context.Context, task *asynq.Task) error {
	var payload PayloadCreateLoginLog
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// create log in DB
	user, err := processor.store.GetUserByEmail(context.Background(), payload.Email)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	_, err = processor.store.CreateLog(context.Background(), db.CreateLogParams{
		UserID:      user.ID,
		Description: fmt.Sprintf("%s logged in at %v", user.Username, time.Now().Format(time.RFC1123)),
	})

	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("user", user.Username).Msg("processed task")
	return nil
}
