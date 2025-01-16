package main

import (
	// Core
	"context"

	"github.com/DMcP89/xlog"

	// All official extensions
	_ "github.com/DMcP89/xlog/extensions/all"
)

func main() {
	xlog.Start(context.Background())
}
