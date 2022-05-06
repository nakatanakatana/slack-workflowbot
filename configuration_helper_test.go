package slackworkflowbot_test

import (
	"testing"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
)

const (
	fooActionID = slackworkflowbot.ActionID("foo-action")
	fooBlockID  = slackworkflowbot.BlockID("foo-block")
	barActionID = slackworkflowbot.ActionID("bar-action")
	barBlockID  = slackworkflowbot.BlockID("bar-block")
)

func createTestBlockActions() slackworkflowbot.BlockActionValues {
	return slackworkflowbot.BlockActionValues{
		string(fooBlockID): {
			string(fooActionID): {
				Value: "foo",
			},
		},
		string(barBlockID): {
			string(barActionID): {
				Value: "bar",
			},
		},
	}
}

func TestCreateTextWorkflowStepInput(t *testing.T) {
	t.Parallel()

	blockActionValues := createTestBlockActions()

	result := slackworkflowbot.CreateTextWorkflowStepInput(blockActionValues, fooActionID, fooBlockID, false)
	got := *result

	if _, ok := got[string(fooBlockID)]; !ok {
		t.Fail()
	}
}

func TestMergeWorkflowStepInput(t *testing.T) {
	t.Parallel()

	blockActionValues := createTestBlockActions()
	resultFoo := slackworkflowbot.CreateTextWorkflowStepInput(blockActionValues, fooActionID, fooBlockID, false)
	resultBar := slackworkflowbot.CreateTextWorkflowStepInput(blockActionValues, barActionID, barBlockID, false)

	result := slackworkflowbot.MergeWorkflowStepInput(resultFoo, resultBar)
	got := *result

	if value, ok := got[string(fooBlockID)]; !ok || value.Value != "foo" {
		t.Fail()
	}

	if value, ok := got[string(barBlockID)]; !ok || value.Value != "bar" {
		t.Fail()
	}
}
