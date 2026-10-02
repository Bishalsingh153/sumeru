package mail

import (
	"context"
	"encoding/json"
	"strings"

	"sumeru/core/applog"
	"sumeru/core/queue"
)

func init() {
	queue.Subscribe("mail", deliverQueuedMail)
}

type mailJob struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	HTMLBody string `json:"html_body,omitempty"`
}

func deliverQueuedMail(ctx context.Context, msg queue.Message) error {
	var job mailJob
	if err := json.Unmarshal(msg.Payload, &job); err != nil {
		return err
	}
	var err error
	if strings.TrimSpace(job.HTMLBody) != "" {
		err = SendMultipart(ctx, job.To, job.Subject, job.Body, job.HTMLBody)
	} else {
		err = Send(ctx, job.To, job.Subject, job.Body)
	}
	if err != nil {
		applog.Warn(ctx, applog.Event{
			Message:   "queued mail delivery failed",
			Component: "mail",
			Operation: "queue_worker",
			Status:    "failed",
			Err:       err,
		})
		return err
	}
	return nil
}

// Enqueue sends mail asynchronously via the in-process queue worker.
func Enqueue(ctx context.Context, to, subject, body string) {
	queue.Publish(ctx, "mail", mailJob{To: to, Subject: subject, Body: body})
}

// EnqueueHTML queues multipart HTML + plain alternative email.
func EnqueueHTML(ctx context.Context, to, subject, htmlBody, textBody string) {
	queue.Publish(ctx, "mail", mailJob{To: to, Subject: subject, Body: textBody, HTMLBody: htmlBody})
}
