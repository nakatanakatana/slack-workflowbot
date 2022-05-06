package slackworkflowbot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/slack-go/slack"
)

//nolint:funlen,cyclop
func CreateInteractionHandler(appCtx AppContext) http.HandlerFunc {
	configureStepCtx := appCtx.configureStep

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

		//nolint:exhaustive
		switch message.Type {
		case slack.InteractionTypeWorkflowStepEdit:
			// https://api.slack.com/workflows/steps#handle_config_view
			configView, ok := configureStepCtx.configView[CallbackID(message.CallbackID)]
			if !ok {
				log.Printf("[WARN] unknown callback id: %s", message.CallbackID)

				return
			}

			err := configView(configureStepCtx, message, "", "")
			if err != nil {
				log.Printf("[ERROR] Failed to open configuration modal in slack: %s", err.Error())
			}

		case slack.InteractionTypeViewSubmission:
			// https://api.slack.com/workflows/steps#handle_view_submission
			saveConfig, ok := configureStepCtx.saveConfig[CallbackID(message.CallbackID)]
			if !ok {
				log.Printf("[WARN] unknown callback id: %s", message.CallbackID)

				return
			}
			//nolint:errcheck
			go saveConfig(configureStepCtx, message)
			w.WriteHeader(http.StatusOK)

		default:
			log.Printf("[WARN] unknown message type: %s", message.Type)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
