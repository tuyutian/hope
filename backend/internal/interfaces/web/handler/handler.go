package handler

import "backend/internal/providers"

type Handlers struct {
}

func InitHandlers(repos *providers.Repositories) *Handlers {
	return &Handlers{}
}
