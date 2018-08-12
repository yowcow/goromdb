package handler

import (
	"fmt"
)

type handlerMap map[string]Handler
type nsHandlerMap map[string]NSHandler

// Multiplexer represents a handler multiplexer
type Multiplexer struct {
	handlers   *handlerMap
	nshandlers *nsHandlerMap
}

// NewMultiplexer creates a Multiplexer
func NewMultiplexer() *Multiplexer {
	handlers := make(handlerMap)
	nsHandlers := make(nsHandlerMap)
	return &Multiplexer{
		handlers:   &handlers,
		nshandlers: &nsHandlers,
	}
}

// RegisterHandler registers a Handler to Multiplexer
func (m *Multiplexer) RegisterHandler(name string, hdr Handler) error {
	if _, ok := (*m.handlers)[name]; ok {
		return fmt.Errorf("handler with name '%s' already registered", name)
	}
	(*m.handlers)[name] = hdr
	return nil
}

// RegisterNSHandler registers a NSHandler to Multiplexer
func (m *Multiplexer) RegisterNSHandler(name string, hdr NSHandler) error {
	if _, ok := (*m.nshandlers)[name]; ok {
		return fmt.Errorf("nshandler with name '%s' already registered", name)
	}
	(*m.nshandlers)[name] = hdr
	return nil
}

// GetHandler returns a Handler with given name
func (m *Multiplexer) GetHandler(name string) (Handler, error) {
	if hdr, ok := (*m.handlers)[name]; ok {
		return hdr, nil
	}
	return nil, fmt.Errorf("handler with name '%s' not registered", name)
}

// GetNSHandler returns a NSHandler with given name
func (m *Multiplexer) GetNSHandler(name string) (NSHandler, error) {
	if hdr, ok := (*m.nshandlers)[name]; ok {
		return hdr, nil
	}
	return nil, fmt.Errorf("nshandler with name '%s' not registered", name)
}
