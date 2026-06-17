package ui

// SpecJSON is the API representation of a spec entry.
type SpecJSON struct {
	ID          string     `json:"id"`
	Dir         string     `json:"dir"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	SpecType    string     `json:"spec_type"`
	Type        string     `json:"type"`
	TaskCount   int        `json:"task_count"`
	Tasks       []TaskJSON `json:"tasks,omitempty"`
}

// TaskJSON is a task summary in list responses.
type TaskJSON struct {
	ID     string `json:"id"`
	File   string `json:"file"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// SpecDetailJSON includes markdown bodies.
type SpecDetailJSON struct {
	SpecJSON
	Markdown string           `json:"markdown"`
	Tasks    []TaskDetailJSON `json:"tasks"`
}

// TaskDetailJSON includes task markdown body.
type TaskDetailJSON struct {
	TaskJSON
	Markdown string `json:"markdown"`
}

// StatusRequest is the body for status PATCH endpoints.
type StatusRequest struct {
	Status string `json:"status"`
}
