package app

import "github.com/docker/cli/cli/compose/types"

// App is the struct for an app, includes:
// - app.yml
// - docker-compose.yml
// - icon.png
type App struct {
	Metadata `yaml:",inline"`

	Type    string        `json:"type" yaml:"type"`
	Compose *types.Config `yaml:"-"`
}

// Metadata is the struct for app.yml
type Metadata struct {
	Title       string   `json:"title" yaml:"title"`
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Categories  []string `json:"categories" yaml:"categories"`
	Platform    string   `json:"platform" yaml:"platform"`
	Note        string   `json:"note" yaml:"note"`
}
