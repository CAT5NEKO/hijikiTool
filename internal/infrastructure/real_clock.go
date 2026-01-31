package infrastructure

import (
	"time"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
)

type RealClock struct{}

func NewRealClock() ports.Clock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}
