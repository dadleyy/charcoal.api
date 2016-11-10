package activity

type Object interface {
	Type() string
	Url() string
}

type Message struct {
	Actor  Object
	Object Object
	Verb   string
}
