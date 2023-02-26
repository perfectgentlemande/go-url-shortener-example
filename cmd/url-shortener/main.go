package main

import (
	"github.com/perfectgentlemande/go-url-shortener-example/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.Module,
	).Run()
}
