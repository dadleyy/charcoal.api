package models

type Object struct {
	Common
}

func (item *Object) Marshal() interface{} {
	return item
}
