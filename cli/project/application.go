package project

import "github.com/docker/cli/cli/compose/types"

// Application is a selfhosted application
type Application struct {
	Name        string
	Description string
	Overview    string
	Category    []string
	Icon        string
	Services    []*types.ServiceConfig
}

// NewApplication create new application
func NewApplication() *Application {
	return &Application{
		Services: []*types.ServiceConfig{},
	}
}
