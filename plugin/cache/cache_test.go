package cache

import "testing"

func TestInit(t *testing.T) {
	Init()
	cache := Cache()
	cache.Do("set", "xx", "sd")
}
