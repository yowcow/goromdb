package storetest

import (
	"fmt"

	"github.com/yowcow/go-romdb/store"
)

type StoreTest struct {
}

func New() store.Store {
	return &StoreTest{}
}

var cacheData = map[string]string{
	"foo": "my test foo",
	"bar": "my test bar!!",
}

func (s StoreTest) Get(key string) (string, error) {
	if v, ok := cacheData[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key '%s' not found", key)
}
