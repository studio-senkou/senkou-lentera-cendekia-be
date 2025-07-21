package main

import "github.com/studio-senkou/lentera-cendekia-be/cmd"

// Senkou Lentera Quiz API
// @title Lentera Quiz API
// @version 1.0
// @description This is the API for the Lentera Quiz application.
// @termsOfService https://example.com/terms
// @contact.name API Support
func main() {
	application := cmd.NewApplication()

	if err := application.Run(); err != nil {
		panic(err)
	}
}

