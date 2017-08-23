package bg

// Object represents a generic record/item that can be represented in a message.
type Object interface {
	Type() string
	URL() string
	Identifier() string
}

// Message is a structure that contains a who, what and "what happened".
type Message struct {
	Actor  Object
	Object Object
	Verb   string
}
