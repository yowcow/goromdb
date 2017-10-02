package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("localhost:11211")
	//mc.Set(&memcache.Item{Key: "foo", Value: []byte("my hoge")})
	//mc.Set(&memcache.Item{Key: "bar", Value: []byte("my bar")})

	{
		it, err := mc.Get("foo")

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(it.Key + " -> " + string(it.Value))
		}
	}

	{
		it, err := mc.Get("hoge")

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(it.Key + " -> " + string(it.Value))
		}
	}

	{
		those, err := mc.GetMulti([]string{"foo", "bar"})

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(those["foo"].Key + " -> " + string(those["foo"].Value))
			fmt.Println(those["bar"].Key + " -> " + string(those["bar"].Value))
		}
	}

	mc.DeleteAll()
}
