package main

type Person struct {
	_id         string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Designation string   `json:"designation,omitempty"`
	Assigment   []string `json:"assigment,omitempty"`
}

type Assignment struct {
	_id    string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Tasks  string `json:"tasks,omitempty"`
	Person string `json:"person,omitempty"`
}
