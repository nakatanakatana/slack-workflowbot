package slackworkflowbot

import "github.com/slack-go/slack"

type (
	StepOutputType string
)

const (
	StepOutputText    = StepOutputType("text")
	StepOutputChannel = StepOutputType("channel")
	StepOutputUser    = StepOutputType("user")
)

type (
	StepOutputConfig struct {
		Name  string
		Type  StepOutputType
		Label string
	}
	StepOutputConfigKey string
	StepOutputs         map[StepOutputConfigKey]StepOutputConfig
)

func CreateOutputsConfig(out StepOutputs) *[]slack.WorkflowStepOutput {
	if out == nil || len(out) == 0 {
		return nil
	}

	result := make([]slack.WorkflowStepOutput, len(out))
	i := 0

	for _, value := range out {
		result[i] = slack.WorkflowStepOutput{
			Name:  value.Name,
			Type:  string(value.Type),
			Label: value.Label,
		}
		i++
	}

	return &result
}

func SetOutputValue(out map[string]string, key StepOutputConfigKey, value string) {
	out[string(key)] = value
}
