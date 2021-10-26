package contentful

import (
	"fmt"
	"net/http"
)

// EnvironmentService model
type EnvironmentService service

// Environment model
type Environment struct {
	Sys *Sys `json:"sys,omitempty"`
}

// GetVersion returns entity version
func (environment *Environment) GetVersion() int {
	version := 1
	if environment.Sys != nil {
		version = environment.Sys.Version
	}

	return version
}

// Get returns a single environment entity
func (service *EnvironmentService) Get(spaceID, environmentID string) (*Environment, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s", spaceID, environmentID)
	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Environment{}, err
	}

	var environment Environment
	if ok := service.c.do(req, &environment); ok != nil {
		return &Environment{}, ok
	}

	return &Environment{}, nil
}