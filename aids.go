package main

import (
	"AIDS/engine"
)

func main() {
	eng := engine.Initialize()
	err := eng.Start()
	if err != nil {
		panic(err)
	}
}
