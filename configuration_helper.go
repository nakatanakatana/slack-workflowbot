package slackworkflowbot

import (
	"github.com/slack-go/slack"
)

type BlockActionValues map[string]map[string]slack.BlockAction

func CreateTextInputBlock(
	aID ActionID,
	bID BlockID,
	name string,
	placeholder string,
	emoji bool,
	verbatim bool,
) *slack.InputBlock {
	text := slack.NewTextBlockObject("plain_text", name, emoji, verbatim)

	var placeholderBlock *slack.TextBlockObject
	if placeholder != "" {
		placeholderBlock = slack.NewTextBlockObject("plain_text", placeholder, emoji, verbatim)
	}

	textElement := slack.NewPlainTextInputBlockElement(placeholderBlock, string(aID))
	textInput := slack.NewInputBlock(string(bID), text, textElement)

	return textInput
}

func CreateTextWorkflowStepInput(
	blockAction BlockActionValues,
	aID ActionID,
	bID BlockID,
	skipVariableReplacement bool,
) *slack.WorkflowStepInputs {
	value := blockAction[string(bID)][string(aID)].Value

	return &slack.WorkflowStepInputs{
		string(bID): slack.WorkflowStepInputElement{
			Value:                   value,
			SkipVariableReplacement: skipVariableReplacement,
		},
	}
}

func MergeWorkflowStepInput(inputs ...*slack.WorkflowStepInputs) *slack.WorkflowStepInputs {
	result := slack.WorkflowStepInputs{}

	for _, tmp := range inputs {
		input := *tmp
		for k, v := range input {
			result[k] = v
		}
	}

	return &result
}
