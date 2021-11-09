package workflow

import (
	"context"
	"fmt"
)

type HelloActivityInput struct {
	Prefix string
	Name   string
}

func HelloActivity(ctx context.Context, input *HelloActivityInput) (string, error) {
	greeting := fmt.Sprintf("Hello: %s!", input.Name)
	return greeting, nil
}
