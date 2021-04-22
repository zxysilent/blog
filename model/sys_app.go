package model

// Goinfo go information
type Goinfo struct {
	ARCH    string `json:"arch"`
	OS      string `json:"os"`
	Version string `json:"version"`
	NumCPU  int    `json:"num_cpu"`
}
