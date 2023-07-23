package main

import (
	"context"
	"go-academy-presentation/pkg/garbage"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(c context.Context) {
	garbage.Migrate()
}

func main() {
	lambda.Start(handler)
}
