package slackworkflowbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

func CreateHandleWorkflowStep(appCtx AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// see: https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			log.Printf("[ERROR] Failed on parsing event: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("my-workflow", eventsAPIEvent, string(body))

		// see: https://api.slack.com/apis/connections/events-api#subscriptions
		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				log.Printf("[ERROR] Failed to decode json message on event url_verification: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
			return
		}

		// see: https://api.slack.com/apis/connections/events-api#receiving_events
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {

			// see: https://api.slack.com/events/workflow_step_execute
			case *slackevents.WorkflowStepExecuteEvent:
				if ev.CallbackID == string(appCtx.workflowStepCallbackID) {
					go appCtx.workflowStep(appCtx, ev.WorkflowStep)
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("[WARN] unknown callbackID: %s", ev.CallbackID)
				return

			default:
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("[WARN] unknown inner event type: %s", eventsAPIEvent.InnerEvent.Type)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[WARN] unknown event type: %s", eventsAPIEvent.Type)
	}
}
