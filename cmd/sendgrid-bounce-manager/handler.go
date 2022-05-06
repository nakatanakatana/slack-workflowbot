package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	bot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	EmailActionID = bot.ActionID("email-input")
	EmailBlockID  = bot.BlockID("email")
	TokenActionID = bot.ActionID("token-input")
	TokenBlockID  = bot.BlockID("token")
)

func saveStepConfig(
	appCtx bot.ConfigureStepContext,
	message slack.InteractionCallback,
) error {
	blockAction := message.View.State.Values
	inEmail := bot.CreateTextWorkflowStepInput(blockAction, EmailActionID, EmailBlockID, false)
	inToken := bot.CreateTextWorkflowStepInput(blockAction, TokenActionID, TokenBlockID, true)

	in := bot.MergeWorkflowStepInput(inEmail, inToken)
	out := &[]slack.WorkflowStepOutput{
		{
			Name:  "message",
			Type:  "text",
			Label: "message",
		},
		{
			Name:  "created_at",
			Type:  "text",
			Label: "created_at",
		},
		{
			Name:  "email",
			Type:  "text",
			Label: "email",
		},
		{
			Name:  "reason",
			Type:  "text",
			Label: "reason",
		},
		{
			Name:  "status",
			Type:  "text",
			Label: "status",
		},
	}

	err := appCtx.SlackClient.SaveWorkflowStepConfiguration(
		message.WorkflowStep.WorkflowStepEditID,
		in,
		out,
	)

	return fmt.Errorf("Slack.SaveWorkflowStepConfiguration Failed: %w", err)
}

func createConfigView() bot.ConfigView {
	return func(
		appCtx bot.ConfigureStepContext,
		message slack.InteractionCallback,
		privateMetaData string,
		externalID string,
	) error {
		emailInput := bot.CreateTextInputBlock(EmailActionID, EmailBlockID, "email", "", false, false)
		tokenInput := bot.CreateTextInputBlock(TokenActionID, TokenBlockID, "token", "", false, false)

		blocks := slack.Blocks{
			BlockSet: []slack.Block{
				emailInput,
				tokenInput,
			},
		}

		cmr := slack.NewConfigurationModalRequest(blocks, privateMetaData, externalID)

		_, err := appCtx.SlackClient.OpenView(message.TriggerID, cmr.ModalViewRequest)
		if err != nil {
			return fmt.Errorf("NewConfigurationModalRequest Failed: %w", err)
		}

		return nil
	}
}

func createCheckStepFunc(sg sendgrid.BounceManager) bot.WorkflowStepFunc {
	return func(
		appCtx bot.StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	) {
		ctx := context.Background()
		inputs := *workflowStep.Inputs

		email, ok := inputs[string(EmailBlockID)]
		if !ok || email.Value == "" {
			message := string(EmailBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		token, ok := inputs[string(TokenBlockID)]
		if !ok || token.Value == "" {
			message := string(TokenBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		out := map[string]string{}

		result, resp, err := sg.GetBounce(email.Value, sendgrid.APIKey(token.Value))

		switch {
		case errors.Is(err, sendgrid.ErrNotFound):
			out["message"] = "not found"
		case err != nil:
			out["message"] = "unexpected error"

			log.Println("get bounce error:", err)
			log.Println("response:", resp)
		default:
			out["message"] = "hit"
			out["created_at"] = result.CreatedAt.Format(time.RFC3339)
			out["email"] = result.Email
			out["reason"] = result.Reason
			out["status"] = result.Status
		}

		err = appCtx.SlackClient.WorkflowStepCompleted(
			ctx,
			workflowStep.WorkflowStepExecuteID,
			&out,
		)
		if err != nil {
			log.Println("workflow step completed error:", err)
		}

		log.Println("Done")
	}
}

func createDeleteStepFunc(sg sendgrid.BounceManager) bot.WorkflowStepFunc {
	return func(
		appCtx bot.StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	) {
		ctx := context.Background()
		inputs := *workflowStep.Inputs

		email, ok := inputs[string(EmailBlockID)]
		if !ok || email.Value == "" {
			message := string(EmailBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		token, ok := inputs[string(TokenBlockID)]
		if !ok || token.Value == "" {
			message := string(TokenBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		out := map[string]string{}

		_, resp, err := sg.DeleteBounce(email.Value, sendgrid.APIKey(token.Value))

		switch {
		case errors.Is(err, sendgrid.ErrNotFound):
			out["message"] = "not found"
		case err != nil:
			out["message"] = "unexpected error"

			log.Println("get bounce error:", err)
			log.Println("response:", resp)
		default:
			out["message"] = "success"
			out["email"] = email.Value
		}
	}
}
