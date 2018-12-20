package watcher

import (
	"context"
)

type Watcher interface {
	Start(context.Context) <-chan string
}
