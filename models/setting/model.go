package setting

// Setting : Project Setting Schema
type Setting struct {
	Streaming Streaming `json:"streaming,omitempty" bson:"streaming,omitempty"`
}

// Streaming : Streaming parametres
type Streaming struct {
	Follow []string `json:"follow" bson:"follow"`
	Track  []string `json:"track" bson:"track"`
}
