package types

// JobStatus represents the state of a rendering job.
type JobStatus struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}
