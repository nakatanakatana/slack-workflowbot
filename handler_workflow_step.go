package slackworkflowbot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

//nolint:funlen
func CreateEventsHandler(appCtx AppContext) http.HandlerFunc {
	stepExecuteCtx := appCtx.stepExecute

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		// see: https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
		body, err := io.ReadAll(r.Body)
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

		switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			// see: https://api.slack.com/apis/connections/events-api#subscriptions
			var r *slackevents.ChallengeResponse

			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Printf("[ERROR] Failed to decode json message on event url_verification: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "text")
			_, _ = w.Write([]byte(r.Challenge))

			return

		case slackevents.CallbackEvent:
			// see: https://api.slack.com/apis/connections/events-api#receiving_events
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			// see: https://api.slack.com/events/workflow_step_execute
			case *slackevents.WorkflowStepExecuteEvent:
				callbackFunc, ok := stepExecuteCtx.workflowStep[CallbackID(ev.CallbackID)]
				if !ok {
					log.Printf("[WARN] unknown callbackID: %s", ev.CallbackID)
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				go callbackFunc(stepExecuteCtx, ev.WorkflowStep)
				w.WriteHeader(http.StatusOK)

				return

			default:
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("[WARN] unknown inner event type: %s", eventsAPIEvent.InnerEvent.Type)

				return
			}
		}
	}
}
