package watcher

import (
	"context"
)

// Watcher defines an interface to a watcher
type Watcher interface {
	Start(context.Context) <-chan string
}
