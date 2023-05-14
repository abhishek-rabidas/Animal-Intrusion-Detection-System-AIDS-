package main

import (
	"AIDS/engine"
	"fmt"
)

func main() {
	fmt.Println("Welcome to AIDS")
	eng := engine.NewAIDSEngine()

	eng.Load()
	eng.Process()
}
