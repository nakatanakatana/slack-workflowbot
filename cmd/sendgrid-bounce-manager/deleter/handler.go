package deleter

import (
	"context"
	"errors"
	"fmt"
	"log"

	bot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	CallbackID = bot.CallbackID("delete-bounce-step")

	EmailBlockID  = bot.BlockID("email")
	EmailActionID = bot.ActionID("email-input")
	TokenBlockID  = bot.BlockID("token")
	TokenActionID = bot.ActionID("token-input")

	EmailInputKey = bot.StepInputConfigKey(EmailBlockID)
	TokenInputKey = bot.StepInputConfigKey(TokenBlockID)

	EmailOutputKey   = bot.StepOutputConfigKey("email")
	MessageOutputKey = bot.StepOutputConfigKey("message")
)

var (
	EmailInputConfig = bot.StepInputConfigText{
		Name:                    string(EmailInputKey),
		Placeholder:             "",
		Multiline:               false,
		Emoji:                   false,
		Verbatim:                false,
		SkipVariableReplacement: false,
	}

	TokenInputConfig = bot.StepInputConfigText{
		Name:                    string(TokenInputKey),
		Placeholder:             "",
		Multiline:               false,
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
	StepOutputConfig = bot.StepOutputs{
		EmailOutputKey:   EmailOutput,
		MessageOutputKey: MessageOutput,
	}
)

func CreateConfigView(cli bot.SlackWorkflowConfigurationClient) bot.ConfigView {
	return func(
		message slack.InteractionCallback,
		privateMetaData string,
		externalID string,
	) error {
		inBlocks := bot.CreateInputsBlock(StepInputConfig, message.WorkflowStep.Inputs)

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

		_, err := cli.OpenView(message.TriggerID, mv)
		if err != nil {
			return fmt.Errorf("NewConfigurationModalRequest Failed: %w", err)
		}

		return nil
	}
}

func CreateSaveStepConfig(cli bot.SlackWorkflowConfigurationClient) bot.SaveConfig {
	return func(
		message slack.InteractionCallback,
	) error {
		blockAction := message.View.State.Values
		in := bot.CreateInputsConfig(blockAction, StepInputConfig)
		out := bot.CreateOutputsConfig(StepOutputConfig)
		err := cli.SaveWorkflowStepConfiguration(
			message.WorkflowStep.WorkflowStepEditID,
			in,
			out,
		)

		return fmt.Errorf("Slack.SaveWorkflowStepConfiguration Failed: %w", err)
	}
}

func CreateStepFunc(cli bot.SlackWorkflowStepExecuteClient, sg sendgrid.BounceManager) bot.WorkflowStep {
	return func(
		workflowStep slackevents.EventWorkflowStep,
	) {
		ctx := context.Background()
		inputs := *workflowStep.Inputs

		email, ok := inputs[string(EmailBlockID)]
		if !ok || email.Value == "" {
			message := string(EmailBlockID) + " is required"
			log.Println("", message)
			_ = cli.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		token, ok := inputs[string(TokenBlockID)]
		if !ok || token.Value == "" {
			message := string(TokenBlockID) + " is required"
			log.Println("", message)
			_ = cli.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		out := map[string]string{}
		bot.SetOutputValue(out, EmailOutputKey, email.Value)

		_, resp, err := sg.DeleteBounce(email.Value, sendgrid.APIKey(token.Value))

		switch {
		case errors.Is(err, sendgrid.ErrNotFound):
			bot.SetOutputValue(out, MessageOutputKey, "not found")

		case err != nil:
			bot.SetOutputValue(out, MessageOutputKey, "unexpected error")

			log.Println("get bounce error:", err)
			log.Println("response:", resp)

		default:
			bot.SetOutputValue(out, MessageOutputKey, "success")
		}

		err = cli.WorkflowStepCompleted(
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
