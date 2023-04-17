package service

import (
	"blog/internal/service/kernel"
	"context"
)

func Init(ctx context.Context) {
	kernel.Init()
	// applet.Init()
}
