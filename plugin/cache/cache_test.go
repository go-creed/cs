package cache

import "testing"

func TestInit(t *testing.T) {
	Init()
	cache := Cache()
	_, err2 := cache.Do("set", "xx", "sd")
	if err2 != nil {
		t.Fatal(err2)
	}
	cache.Do("del", "xx")

	cache.Do("setex", "jwt", "100", "jsfla")
}
