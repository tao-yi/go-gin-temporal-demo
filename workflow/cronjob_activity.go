package workflow

import (
	"context"
	"fmt"
	"time"
)

func CrobJobActivity(ctx context.Context) error {
	fmt.Println(time.Now().UTC())
	return nil
}
