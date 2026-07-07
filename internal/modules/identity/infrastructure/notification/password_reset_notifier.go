// =========================================================
// File: internal/modules/identity/infrastructure/notification/password_reset_notifier.go.go
// =========================================================

package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ruangwali/internal/modules/identity/application/ports"
)

type PasswordResetNotifier struct {
	client   *http.Client
	endpoint string
	apiKey   string
}

func NewPasswordResetNotifier(
	client *http.Client,
	endpoint string,
	apiKey string,
) (*PasswordResetNotifier, error) {
	if client == nil {
		return nil, errors.New(
			"password reset notifier: http client nil",
		)
	}

	endpoint = strings.TrimSpace(
		endpoint,
	)

	apiKey = strings.TrimSpace(
		apiKey,
	)

	if endpoint == "" {
		return nil, errors.New(
			"password reset notifier: endpoint wajib diisi",
		)
	}

	if apiKey == "" {
		return nil, errors.New(
			"password reset notifier: api key wajib diisi",
		)
	}

	return &PasswordResetNotifier{
		client:   client,
		endpoint: endpoint,
		apiKey:   apiKey,
	}, nil
}

func (n *PasswordResetNotifier) SendPasswordReset(
	ctx context.Context,
	notification ports.PasswordResetNotification,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	payload := passwordResetPayload{
		UserID: notification.UserID.String(),
		Email:  notification.Email.String(),
		Token:  notification.Token,

		ExpiresAt: notification.ExpiresAt.UTC(),
	}

	body, err := json.Marshal(
		payload,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat payload password reset: %w",
			err,
		)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		n.endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf(
			"gagal membuat request password reset: %w",
			err,
		)
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	request.Header.Set(
		"X-API-Key",
		n.apiKey,
	)

	response, err := n.client.Do(
		request,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal mengirim password reset notification: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusOK &&
		response.StatusCode < http.StatusMultipleChoices {
		_, _ = io.Copy(
			io.Discard,
			response.Body,
		)

		return nil
	}

	responseBody, readErr := io.ReadAll(
		io.LimitReader(
			response.Body,
			4096,
		),
	)
	if readErr != nil {
		return fmt.Errorf(
			"password reset notifier menerima status %d",
			response.StatusCode,
		)
	}

	message := strings.TrimSpace(
		string(responseBody),
	)

	if message == "" {
		return fmt.Errorf(
			"password reset notifier menerima status %d",
			response.StatusCode,
		)
	}

	return fmt.Errorf(
		"password reset notifier menerima status %d: %s",
		response.StatusCode,
		message,
	)
}

type passwordResetPayload struct {
	UserID string `json:"user_id"`

	Email string `json:"email"`

	Token string `json:"token"`

	ExpiresAt time.Time `json:"expires_at"`
}

var _ ports.PasswordResetNotifier = (*PasswordResetNotifier)(nil)
