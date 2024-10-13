package whatsapp

// Controls Cache and Events from current WhatsApp service
type IWhatsappHandlers interface {

	// Process a single message
	Message(*WhatsappMessage)

	// Update message status information
	MessageStatusUpdate(id string, status WhatsappMessageStatus) bool

	// Update read receipt status
	Receipt(*WhatsappMessage)

	// Event
	LoggedOut(string)

	GetLeading() *WhatsappMessage

	GetById(id string) (*WhatsappMessage, error)

	OnConnected()

	OnDisconnected()

	IsInterfaceNil() bool
}
