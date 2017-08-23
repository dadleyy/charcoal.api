package models

// Activity records represent an action, a thing and the person who took the action on the thing.
type Activity struct {
	Common
	Type       string `json:"verb"`
	ActorType  string `json:"actor_type"`
	ActorURL   string `json:"actor_url"`
	ActorUUID  string `json:"actor_uuid"`
	ObjectType string `json:"object_type"`
	ObjectURL  string `json:"object_url"`
	ObjectUUID string `json:"object_uuid"`
}

// TableName is used by GORM to prevent auto-pluralization.
func (activity Activity) TableName() string {
	return "activity"
}
