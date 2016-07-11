package common

// MessageType designate the type of a message send from a client 
type MessageType int
const (
    // FUNC for functionnal messages ie technical messages from the client to the server
	FUNC MessageType = iota
    // CLASSIQUE message for messages sent by the end user
	CLASSIQUE
)


// ConnectionStatus is self explained
type ConnectionStatus int
const (
	JOINING ConnectionStatus = iota
	LEAVING
)
