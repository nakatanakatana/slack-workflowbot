package slackworkflowbot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/slack-go/slack"
)

func CreateHandleInteraction(appCtx AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonStr, err := url.QueryUnescape(string(body)[8:])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var message slack.InteractionCallback
		if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
			log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("interaction", message.Type, string(body))

		switch message.Type {
		case slack.InteractionTypeWorkflowStepEdit:
			// https://api.slack.com/workflows/steps#handle_config_view
			err := appCtx.replyWithConfigurationView(appCtx, message, "", "")
			if err != nil {
				log.Printf("[ERROR] Failed to open configuration modal in slack: %s", err.Error())
			}

		case slack.InteractionTypeViewSubmission:
			// https://api.slack.com/workflows/steps#handle_view_submission
			go appCtx.saveUserSettingsForWorkflowStep(appCtx, message)
			w.WriteHeader(http.StatusOK)

		default:
			log.Printf("[WARN] unknown message type: %s", message.Type)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
