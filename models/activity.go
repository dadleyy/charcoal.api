package models

type Activity struct {
	Common
	Type       string `json:"verb"`
	ActorType  string `json:"actor_type"`
	ActorUrl   string `json:"actor_url"`
	ActorUuid  string `json:"actor_uuid"`
	ObjectType string `json:"object_type"`
	ObjectUrl  string `json:"object_url"`
	ObjectUuid string `json:"object_uuid"`
}

func (activity Activity) TableName() string {
	return "activity"
}

func (activity Activity) Public() interface{} {
	return activity
}
