package slackworkflowbot

import (
	"github.com/slack-go/slack"
)

type (
	StepInputConfig struct {
		ActionID ActionID
		BlockID  BlockID
		Config   StepInputConfigDetail
	}
	StepInputConfigKey string

	StepInputConfigDetail interface {
		Type() slack.InputType
	}

	StepInputs map[StepInputConfigKey]StepInputConfig
)

type StepInputConfigText struct {
	Name                    string
	Placeholder             string
	Emoji                   bool
	Multiline               bool
	Verbatim                bool
	SkipVariableReplacement bool
}

func (c StepInputConfigText) Type() slack.InputType {
	return slack.InputTypeText
}

type StepInputConfigCheckbox struct {
	Description             string
	Options                 []string
	Emoji                   bool
	Verbatim                bool
	SkipVariableReplacement bool
}

func (c StepInputConfigCheckbox) Type() slack.InputType {
	return slack.InputTypeText
}

func CreateInputsConfig(blockAction BlockActionValues, inputs StepInputs) *slack.WorkflowStepInputs {
	in := make([]*slack.WorkflowStepInputs, len(inputs))
	i := 0

	for _, value := range inputs {
		switch cfg := value.Config.(type) {
		case StepInputConfigText:
			in[i] = CreateTextWorkflowStepInput(
				blockAction,
				value.ActionID,
				value.BlockID,
				cfg.SkipVariableReplacement,
			)
			i++

		case StepInputConfigCheckbox:
			in[i] = CreateTextWorkflowStepInput(
				blockAction,
				value.ActionID,
				value.BlockID,
				cfg.SkipVariableReplacement,
			)
			i++

		default:
			continue
		}
	}

	return MergeWorkflowStepInput(in[:i]...)
}

func CreateInputsBlock(inputs StepInputs, values *slack.WorkflowStepInputs) map[StepInputConfigKey]*slack.InputBlock {
	results := make(map[StepInputConfigKey]*slack.InputBlock)
	inputValues := *values

	for key, value := range inputs {
		input := ""
		if v, ok := inputValues[string(key)]; ok {
			input = v.Value
		}

		switch cfg := value.Config.(type) {
		case StepInputConfigText:
			results[key] = CreateTextInputBlock(
				value.ActionID,
				value.BlockID,
				cfg.Name,
				cfg.Placeholder,
				input,
				cfg.Multiline,
				cfg.Emoji,
				cfg.Verbatim,
			)

		case StepInputConfigCheckbox:
			results[key] = CreateCheckboxBlock(
				value.ActionID,
				value.BlockID,
				cfg.Description,
				cfg.Options,
				cfg.Emoji,
				cfg.Verbatim,
			)

		default:
			continue
		}
	}

	return results
}
