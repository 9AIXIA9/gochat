package kafka

import (
	"context"
	"fmt"
	"gochat/internal/authorization/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedHandler(publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		fmt.Println("user created event received")
		userCreatedEvent, err := domain.ToUserCreatedEvent(e)
		if err != nil {
			return err
		}

		emailSendingRequestedEvent, err := notificationDomain.NewEmailSendingRequestedEvent(
			userCreatedEvent.ID(),
			kernel.UserID(userCreatedEvent.AggregateID()),
			userCreatedEvent.Email(),
			"signup",
			"Welcome to GoChat!",
			signupEmailTemplate,
		)
		if err != nil {
			return err
		}

		return publisher.Publish([]event.Event{emailSendingRequestedEvent})
	}
}

const signupEmailTemplate = `
<h1>Welcome to GoChat!</h1>
<p>We're thrilled to have you join our community. Now you're all set to start seamless and secure conversations.</p>
<p>To get the most out of your GoChat experience, you can:</p>
<ul>
  <li><strong>Complete your profile:</strong> Let people know who you are.</li>
  <li><strong>Find friends:</strong> Connect with your contacts or explore public groups.</li>
  <li><strong>Start chatting:</strong> Send your first message – it's that easy!</li>
</ul>
<p>If you have any questions, feel free to reply to this email. We're here to help!</p>
<p>Cheers,<br>The GoChat Team</p>
`
