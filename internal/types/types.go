// Package types contains shared type definitions for Skooma
package types

// Template represents a project template
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	RepoURL     string `json:"repo_url"`
}

// Config represents the Skooma configuration
type Config struct {
	Templates map[string]Template `json:"templates"`
}

// ProjectData holds the data collected from the user to populate the project templates.
type ProjectData struct {
	Name     string   `json:"name"`
	RootDir  string   `json:"root_dir"`
	Template Template `json:"template"`
	RepoURL  string   `json:"repo_url"`
	Author   string   `json:"author"`
}
