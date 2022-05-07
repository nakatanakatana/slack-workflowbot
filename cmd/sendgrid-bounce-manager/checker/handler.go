package checker

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
	CallbackID = bot.CallbackID("check-bounce-step")

	EmailBlockID  = bot.BlockID("email")
	EmailActionID = bot.ActionID("email-input")
	TokenBlockID  = bot.BlockID("token")
	TokenActionID = bot.ActionID("token-input")

	EmailInputKey = bot.StepInputConfigKey(EmailBlockID)
	TokenInputKey = bot.StepInputConfigKey(TokenBlockID)

	EmailOutputKey     = bot.StepOutputConfigKey("email")
	MessageOutputKey   = bot.StepOutputConfigKey("message")
	CreatedAtOutputKey = bot.StepOutputConfigKey("created_at")
	ReasonOutputKey    = bot.StepOutputConfigKey("reason")
	StatusOutputKey    = bot.StepOutputConfigKey("status")
)

var (
	EmailInputConfig = bot.StepInputConfigText{
		Name:                    string(EmailInputKey),
		Placeholder:             "",
		Emoji:                   false,
		Verbatim:                false,
		SkipVariableReplacement: false,
	}

	TokenInputConfig = bot.StepInputConfigText{
		Name:                    string(TokenInputKey),
		Placeholder:             "",
		Emoji:                   false,
		Verbatim:                false,
		SkipVariableReplacement: true,
	}
	StepInputConfig = bot.StepInputs{
		EmailInputKey: bot.StepInputConfig{
			ActionID: EmailActionID,
			BlockID:  EmailBlockID,
			Config:   EmailInputConfig,
		},
		TokenInputKey: bot.StepInputConfig{
			ActionID: TokenActionID,
			BlockID:  TokenBlockID,
			Config:   TokenInputConfig,
		},
	}
)

var (
	EmailOutput = bot.StepOutputConfig{
		Name:  string(EmailOutputKey),
		Type:  bot.StepOutputText,
		Label: string(EmailOutputKey),
	}

	MessageOutput = bot.StepOutputConfig{
		Name:  string(MessageOutputKey),
		Type:  bot.StepOutputText,
		Label: string(MessageOutputKey),
	}

	CreatedAtOutput = bot.StepOutputConfig{
		Name:  string(CreatedAtOutputKey),
		Type:  bot.StepOutputText,
		Label: string(CreatedAtOutputKey),
	}

	ReasonOutput = bot.StepOutputConfig{
		Name:  string(ReasonOutputKey),
		Type:  bot.StepOutputText,
		Label: string(ReasonOutputKey),
	}

	StatusOutput = bot.StepOutputConfig{
		Name:  string(StatusOutputKey),
		Type:  bot.StepOutputText,
		Label: string(StatusOutputKey),
	}

	StepOutputConfig = bot.StepOutputs{
		EmailOutputKey:     EmailOutput,
		MessageOutputKey:   MessageOutput,
		CreatedAtOutputKey: CreatedAtOutput,
		ReasonOutputKey:    ReasonOutput,
		StatusOutputKey:    StatusOutput,
	}
)

func CreateConfigView() bot.ConfigView {
	return func(
		appCtx bot.ConfigureStepContext,
		message slack.InteractionCallback,
		privateMetaData string,
		externalID string,
	) error {
		inBlocks := bot.CreateInputsBlock(StepInputConfig)

		blocks := slack.Blocks{
			BlockSet: []slack.Block{
				inBlocks[EmailInputKey],
				inBlocks[TokenInputKey],
			},
		}

		mv := slack.ModalViewRequest{
			Type:            slack.VTWorkflowStep,
			Title:           nil, // slack configuration modal must not have a title!
			Blocks:          blocks,
			CallbackID:      message.CallbackID,
			PrivateMetadata: privateMetaData,
			ExternalID:      externalID,
		}

		_, err := appCtx.SlackClient.OpenView(message.TriggerID, mv)
		if err != nil {
			return fmt.Errorf("NewConfigurationModalRequest Failed: %w", err)
		}

		return nil
	}
}

func SaveStepConfig(
	appCtx bot.ConfigureStepContext,
	message slack.InteractionCallback,
) error {
	blockAction := message.View.State.Values
	in := bot.CreateInputsConfig(blockAction, StepInputConfig)
	out := bot.CreateOutputsConfig(StepOutputConfig)

	err := appCtx.SlackClient.SaveWorkflowStepConfiguration(
		message.WorkflowStep.WorkflowStepEditID,
		in,
		out,
	)

	return fmt.Errorf("Slack.SaveWorkflowStepConfiguration Failed: %w", err)
}

func CreateStepFunc(sg sendgrid.BounceManager) bot.WorkflowStepFunc {
	return func(
		appCtx bot.StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	) {
		ctx := context.Background()
		inputs := *workflowStep.Inputs

		email, ok := inputs[string(EmailBlockID)]
		if !ok || email.Value == "" {
			message := string(EmailBlockID) + " is required"
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		token, ok := inputs[string(TokenBlockID)]
		if !ok || token.Value == "" {
			message := string(TokenBlockID) + " is required"
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		out := map[string]string{}
		bot.SetOutputValue(out, EmailOutputKey, email.Value)
		result, resp, err := sg.GetBounce(email.Value, sendgrid.APIKey(token.Value))

		switch {
		case errors.Is(err, sendgrid.ErrNotFound):
			bot.SetOutputValue(out, MessageOutputKey, "not found")

		case err != nil:
			bot.SetOutputValue(out, MessageOutputKey, "unexpected error")

			log.Println("get bounce error:", err)
			log.Println("response:", resp)

		default:
			bot.SetOutputValue(out, MessageOutputKey, "matched")
			bot.SetOutputValue(out, CreatedAtOutputKey, result.CreatedAt.Format(time.RFC3339))
			bot.SetOutputValue(out, ReasonOutputKey, result.Reason)
			bot.SetOutputValue(out, StatusOutputKey, result.Status)
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