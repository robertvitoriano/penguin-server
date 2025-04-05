package receiveEvents

type ReceiveEvent string

const (
	StartGame ReceiveEvent = "start_game"
	Close     ReceiveEvent = "close"
	// Add other ReceiveEvents here as needed
)
