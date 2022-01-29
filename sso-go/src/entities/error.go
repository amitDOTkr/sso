package entities

type Error struct {
	Type   string `json:"type,omitempty" bson:"type,omitempty"`
	Detail string `json:"detail,omitempty" bson:"detail,omitempty"`
}
