package main

import "im/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}
