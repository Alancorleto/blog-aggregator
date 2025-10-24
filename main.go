package main

import (
	"fmt"

	config "github.com/alancorleto/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	err = cfg.SetUser("alan")
	if err != nil {
		fmt.Println("Error setting user:", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println("Current User Name:", cfg.CurrentUserName)
}
