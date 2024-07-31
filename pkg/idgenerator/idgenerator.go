package idgenerator

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"
)

type Generator interface {
	Generate() string
}

type idGenerator struct {
	entropy *ulid.MonotonicEntropy
}

func New() Generator {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return &idGenerator{entropy: entropy}
}

func (generator *idGenerator) Generate() string {
	return strings.ToLower(ulid.MustNew(ulid.Timestamp(time.Now()), generator.entropy).String())
}
