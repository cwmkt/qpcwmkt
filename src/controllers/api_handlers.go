package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	models "github.com/nocodeleaks/quepasa/models"
)

const CurrentAPIVersion string = "v4"

func RegisterAPIControllers(r chi.Router) {

	aliases := []string{"/current", "", "/" + CurrentAPIVersion}
	for _, endpoint := range aliases {

		// CONTROL METHODS ************************
		// ----------------------------------------
		r.Get(endpoint+"/info", InformationController)
		r.Patch(endpoint+"/info", InformationController)
		r.Delete(endpoint+"/info", InformationController)

		r.Get(endpoint+"/scan", ScannerController)
		r.Get(endpoint+"/paircode", PairCodeController)

		r.Get(endpoint+"/command", CommandController)

		// ----------------------------------------
		// CONTROL METHODS ************************

		// SENDING MSG ----------------------------
		// ----------------------------------------

		r.Get(endpoint+"/message/{messageid}", GetMessageController)
		r.Get(endpoint+"/message", GetMessageController)

		r.Delete(endpoint+"/message/{messageid}", RevokeController)
		r.Delete(endpoint+"/message", RevokeController)

		// used to send alert msgs via url, triggers on monitor systems like zabbix
		r.Get(endpoint+"/send", SendAny)

		r.Post(endpoint+"/send", SendAny)
		r.Post(endpoint+"/send/{chatid}", SendAny)
		r.Post(endpoint+"/sendtext", SendText)
		r.Post(endpoint+"/sendtext/{chatid}", SendText)

		// SENDING MSG ATTACH ---------------------

		// deprecated, discard/remove on next version
		r.Post(endpoint+"/senddocument", SendDocumentAPIHandlerV2)

		r.Post(endpoint+"/sendurl", SendAny)
		r.Post(endpoint+"/sendbinary/{chatid}/{filename}/{text}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary/{chatid}/{filename}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary/{chatid}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary", SendDocumentFromBinary)
		r.Post(endpoint+"/sendencoded", SendAny)

		// ----------------------------------------
		// SENDING MSG ----------------------------

		r.Get(endpoint+"/receive", ReceiveAPIHandler)
		r.Post(endpoint+"/attachment", AttachmentAPIHandlerV2)

		r.Get(endpoint+"/download/{messageid}", DownloadController)
		r.Get(endpoint+"/download", DownloadController)

		// PICTURE INFO | DATA --------------------
		// ----------------------------------------

		r.Post(endpoint+"/picinfo", PictureController)
		r.Get(endpoint+"/picinfo/{chatid}/{pictureid}", PictureController)
		r.Get(endpoint+"/picinfo/{chatid}", PictureController)
		r.Get(endpoint+"/picinfo", PictureController)

		r.Post(endpoint+"/picdata", PictureController)
		r.Get(endpoint+"/picdata/{chatid}/{pictureid}", PictureController)
		r.Get(endpoint+"/picdata/{chatid}", PictureController)
		r.Get(endpoint+"/picdata", PictureController)

		// ----------------------------------------
		// PICTURE INFO | DATA --------------------

		r.Post(endpoint+"/webhook", WebhookController)
		r.Get(endpoint+"/webhook", WebhookController)
		r.Delete(endpoint+"/webhook", WebhookController)

		// INVITE METHODS ************************
		// ----------------------------------------

		r.Get(endpoint+"/invite", InviteController)
		r.Get(endpoint+"/invite/{chatid}", InviteController)

		// ----------------------------------------
		// INVITE METHODS ************************

		r.Get(endpoint+"/contacts", ContactsController)
		r.Post(endpoint+"/isonwhatsapp", IsOnWhatsappController)

		// IF YOU LOVE YOUR FREEDOM, DO NOT USE THAT
		// IT WAS DEVELOPED IN A MOMENT OF WEAKNESS
		// DONT BE THAT GUY !
		r.Post(endpoint+"/spam", Spam)
	}
}

func CommandController(w http.ResponseWriter, r *http.Request) {
	// setting default response type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpResponse{}

	server, err := GetServer(r)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	action := models.GetRequestParameter(r, "action")
	switch action {
	case "start":
		err = server.Start()
		if err == nil {
			response.ParseSuccess("started")
		}
	case "stop":
		err = server.Stop("command")
		if err == nil {
			response.ParseSuccess("stopped")
		}
	case "restart":
		err = server.Restart()
		if err == nil {
			response.ParseSuccess("restarted")
		}
	case "status":
		status := server.GetStatus()
		response.ParseSuccess(status.String())
	case "groups":
		err := models.ToggleGroups(server)
		if err == nil {
			message := "groups toggled: " + server.Groups.String()
			response.ParseSuccess(message)
		}
	case "broadcasts":
		err := models.ToggleBroadcasts(server)
		if err == nil {
			message := "broadcasts toggled: " + server.Broadcasts.String()
			response.ParseSuccess(message)
		}
	case "readreceipts":
		err := models.ToggleReadReceipts(server)
		if err == nil {
			message := "readreceipts toggled: " + server.ReadReceipts.String()
			response.ParseSuccess(message)
		}
	case "calls":
		err := models.ToggleCalls(server)
		if err == nil {
			message := "calls toggled: " + server.Calls.String()
			response.ParseSuccess(message)
		}
	default:
		err = fmt.Errorf("invalid action: {%s}, try {start,stop,restart,status,groups}", action)
	}

	if err != nil {
		response.ParseError(err)
	}

	RespondInterface(w, response)
}
