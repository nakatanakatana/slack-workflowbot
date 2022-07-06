package slackworkflowbot

import (
	"strings"

	"github.com/slack-go/slack"
)

type BlockActionValues map[string]map[string]slack.BlockAction

func createOptionBlockObjects(options []string, emoji, verbatim bool) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))

	for _, o := range options {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, o, emoji, verbatim)
		optionTextBlock := slack.NewOptionBlockObject(o, optionText, nil)
		optionBlockObjects = append(optionBlockObjects, optionTextBlock)
	}

	return optionBlockObjects
}

// TextInput

func CreateTextInputBlock(
	aID ActionID,
	bID BlockID,
	name string,
	placeholder string,
	multiline bool,
	emoji bool,
	verbatim bool,
) *slack.InputBlock {
	text := slack.NewTextBlockObject(slack.PlainTextType, name, emoji, verbatim)

	var placeholderBlock *slack.TextBlockObject
	if placeholder != "" {
		placeholderBlock = slack.NewTextBlockObject(slack.PlainTextType, placeholder, emoji, verbatim)
	}

	textElement := slack.NewPlainTextInputBlockElement(placeholderBlock, string(aID))
	textElement.Multiline = multiline
	textInput := slack.NewInputBlock(string(bID), text, nil, textElement)

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

// Checkbox

func CreateCheckboxBlock(
	aID ActionID,
	bID BlockID,
	description string,
	options []string,
	emoji bool,
	verbatim bool,
) *slack.InputBlock {
	descriptionBlock := slack.NewTextBlockObject(slack.PlainTextType, description, emoji, verbatim)
	checkboxOptions := createOptionBlockObjects(options, emoji, verbatim)

	checkboxOptionsBlock := slack.NewCheckboxGroupsBlockElement(string(aID), checkboxOptions...)

	return slack.NewInputBlock(string(bID), descriptionBlock, nil, checkboxOptionsBlock)
}

func CreateCheckboxWorkflowStepInput(
	blockAction BlockActionValues,
	aID ActionID,
	bID BlockID,
	skipVariableReplacement bool,
) *slack.WorkflowStepInputs {
	options := blockAction[string(bID)][string(aID)].SelectedOptions
	values := make([]string, len(options))

	for i, o := range options {
		values[i] = o.Value
	}

	return &slack.WorkflowStepInputs{
		string(bID): slack.WorkflowStepInputElement{
			Value:                   strings.Join(values, ","),
			SkipVariableReplacement: skipVariableReplacement,
		},
	}
}

// helper

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
