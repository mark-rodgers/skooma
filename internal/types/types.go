// Package types contains shared type definitions for Skooma
package types

// Template represents a project template
type Template struct {
	Description string `json:"description"`
	Repo        string `json:"repo"`
	Author      string `json:"author"`
}

// Config represents the Skooma configuration
type Config struct {
	Templates map[string]Template `json:"templates"`
}