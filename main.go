package main

import (
	"net/http"
)

func main() {
	QueueTest()
}

// Controller houses all the methods we need for the HTTP handlers
type Controller struct {
	CustomerService
	TicketService
}

// TicketRecord contains a single conversation between customer and service desk representative
type TicketRecord struct {
	ID            int
	ProjectID     int
	TaskID        int
	CustomerEmail string
	Closed        bool
}

// CustomerService handles the interface between this system and the customer
type CustomerService interface {
	Send(customerEmail string, subject string, body string)
}

// TicketService handles the interface between this system and the ticketing system
type TicketService interface {
	TaskCreate(projectID int, body string) int
	TaskClose(taskID int)
	Send(projectID int, taskID int, body string)
}

// EmailReceiveHandler will handle emails received from SendGrid or equivalent
// The email will either create a new task, or append a comment to an existing task
func (c *Controller) EmailReceiveHandler(w http.ResponseWriter, r *http.Request) {
	// email := c.EmailReceive(r.Body)
	// c.CommentSend(email.ProjectID)
}

// WebhookHandler will handle webhooks received from Kanboard or equivalent
// The comment is a webhook resulting from a new comment on an existing task
// The response will be forwarded onto the customer
// This handler will also handle closing of a task
func (c *Controller) WebhookHandler(w http.ResponseWriter, r *http.Request) {

}
