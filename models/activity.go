package models

type Activity struct {
	Common
	Type string `json:"verb"`
	ActorType string `json:"actor_type"`
	ActorUrl string `json:"actor_url"`
	ObjectType string `json:"object_type"`
	ObjectUrl string `json:"object_url"`
}

func (activity Activity) TableName() string {
	return "activity"
}

func (activity Activity) Public() interface{} {
	return activity
}
