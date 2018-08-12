package handler

import (
	"fmt"
)

type handlerMap map[string]Handler
type nsHandlerMap map[string]NSHandler

type Multiplexer struct {
	handlers   *handlerMap
	nshandlers *nsHandlerMap
}

func NewMultiplexer() *Multiplexer {
	handlers := make(handlerMap)
	nsHandlers := make(nsHandlerMap)
	return &Multiplexer{
		handlers:   &handlers,
		nshandlers: &nsHandlers,
	}
}

func (m *Multiplexer) RegisterHandler(name string, hdr Handler) error {
	if _, ok := (*m.handlers)[name]; ok {
		return fmt.Errorf("handler with name '%s' already registered", name)
	}
	(*m.handlers)[name] = hdr
	return nil
}

func (m *Multiplexer) RegisterNSHandler(name string, hdr NSHandler) error {
	if _, ok := (*m.nshandlers)[name]; ok {
		return fmt.Errorf("nshandler with name '%s' already registered", name)
	}
	(*m.nshandlers)[name] = hdr
	return nil
}

func (m *Multiplexer) GetHandler(name string) (Handler, error) {
	if hdr, ok := (*m.handlers)[name]; ok {
		return hdr, nil
	}
	return nil, fmt.Errorf("handler with name '%s' not registered", name)
}

func (m *Multiplexer) GetNSHandler(name string) (NSHandler, error) {
	if hdr, ok := (*m.nshandlers)[name]; ok {
		return hdr, nil
	}
	return nil, fmt.Errorf("nshandler with name '%s' not registered", name)
}
