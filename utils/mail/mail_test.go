package mail_test

import (
	"testing"

	mail "github.com/studio-senkou/lentera-cendekia-be/utils/mail"
)

func TestSendEmail(t *testing.T) {
	// This is a placeholder for the actual test implementation.
	// You would typically call the function that sends an email and check if it behaves as expected.
	email := mail.NewMail(
		"ajhmdni02@gmail.com",
		"Test Email",
		"This is a test email body.",
	)

	email.Send()	
}
