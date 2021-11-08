package activity

import "fmt"

func HelloActivity(name string) (string, error) {
	greeting := fmt.Sprintf("Hello: %s!", name)
	return greeting, nil
}
