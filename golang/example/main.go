package main

import (
	"fmt"
	"os/user"
)

func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	fmt.Println(greet("World"))

	u, err := user.Current()
	if err != nil {
		fmt.Printf("could not get current user: %v\n", err)
		return
	}
	fmt.Printf("Running as user: %s\n", u.Username)
}
