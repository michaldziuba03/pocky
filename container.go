package main

import "github.com/google/uuid"

type Container struct {
	ID     string
	SID    string
	limits CGroup
}

func NewContainer() *Container {
	id := uuid.New().String()
	sid := id[:8] // for display only

	container := &Container{
		ID:     id,
		SID:    sid,
		limits: CGroup{},
	}

	container.limits.InitCGroup(container)
	return container
}
