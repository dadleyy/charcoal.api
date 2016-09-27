package models

type Activity struct {
	Common
}

func (item *Activity) Marshal() interface{} {
	return item
}
