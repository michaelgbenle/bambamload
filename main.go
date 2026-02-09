package main

import (
	b "bambamload/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	b.Start()
}
