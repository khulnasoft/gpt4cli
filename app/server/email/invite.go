package email

import (
	"fmt"
	"os"

	"github.com/gen2brain/beeep"
)

func SendInviteEmail(email, inviteeFirstName, inviterName, orgName string) error {
	// Check if the environment is production
	if os.Getenv("GOENV") == "production" {
		// Production environment - send email using AWS SES
		subject := fmt.Sprintf("%s, you've been invited to join %s on Gpt4cli", inviteeFirstName, orgName)

		htmlBody := fmt.Sprintf(`<p>Hi %s,</p><p>%s has invited you to join the org <strong>%s</strong> on <a href="https://gpt4cli.khulnasoft.com">Gpt4cli.</a></p><p>Gpt4cli is a terminal-based AI programming engine for complex tasks.</p><p>To accept the invite, first <a href="https://docs-gpt4cli.khulnasoft.com/install/">install Gpt4cli</a>, then open a terminal and run 'gpt4cli sign-in'. Enter '%s' when asked for your email and follow the prompts from there.</p><p>If you have questions, feedback, or run into a problem, you can reply directly to this email, <a href="https://github.com/khulnasoft/gpt4cli/discussions">start a discussion</a>, or <a href="https://github.com/khulnasoft/gpt4cli/issues">open an issue.</a></p>`, inviteeFirstName, inviterName, orgName, email)

		textBody := fmt.Sprintf(`Hi %s,\n\n%s has invited you to join the org %s on Gpt4cli.\n\nGpt4cli is a terminal-based AI programming engine for complex tasks.\n\nTo accept the invite, first install Gpt4cli (https://docs-gpt4cli.khulnasoft.com/install/), then open a terminal and run 'gpt4cli sign-in'. Enter '%s' when asked for your email and follow the prompts from there.\n\nIf you have questions, feedback, or run into a problem, you can reply directly to this email, start a discussion (https://github.com/khulnasoft/gpt4cli/discussions), or open an issue (https://github.com/khulnasoft/gpt4cli/issues).`, inviteeFirstName, inviterName, orgName, email)

		if os.Getenv("IS_CLOUD") == "" {
			return sendEmailViaSMTP(email, subject, htmlBody, textBody)
		} else {
			return SendEmailViaSES(email, subject, htmlBody, textBody)
		}
	} else {
		// Send notification
		err := beeep.Notify("Invite Sent", fmt.Sprintf("Invite sent to %s (email not sent in development)", email), "")
		if err != nil {
			return fmt.Errorf("error sending notification in dev: %v", err)
		}
	}

	return nil
}
