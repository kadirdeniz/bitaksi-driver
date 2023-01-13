package main

import (
	_ "driver/docs"
	"driver/tools/fiber"
	"fmt"
)

func main() {
	fmt.Println("Driver Service")

	fiber.Router(8001)
}
