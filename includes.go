package contentful

// Includes model
type Includes struct {
	Asset []*Asset `json:"Asset,omitempty"`
	Entry []*Entry `json:"Entry,omitempty"`
}