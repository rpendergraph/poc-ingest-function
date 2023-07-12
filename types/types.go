package types

type EventSink[T any] interface {
	PublishEvents([]T) error
}

type GovernedEvent struct {
	ID             string   `json:"id"`
	ContainsPII    bool     `json:"containsPII"`
	ContainsPHI    bool     `json:"containsPHI"`
	LineOfBusiness string   `json:"lineOfBusiness"`
	Source         string   `json:"source"`
	Destination    string   `json:"destination"`
	SizeBytes      uint     `json:"sizeBytes"`
	Tags           []string `json:"tags"`
}
