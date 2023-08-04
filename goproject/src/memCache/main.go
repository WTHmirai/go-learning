package main

import (
	"fmt"
	"memCache/cache"
)

func main() {
	tempcache := cache.NewCache()
	_ = tempcache.SetMaxMemory("10GB")

	fmt.Println(cache.ParseOcSize("ahfshasd"))
	fmt.Println(cache.ParseOcSize(map[string]string{
		"闸门萨达": "ashdasfj",
	}))
}
