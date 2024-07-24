package otp

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator interface {
	Generate() string
}

type generator struct {
	randSource *rand.Rand
}

func New() Generator {
	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)
	return &generator{randSource: randSource}
}

func (o *generator) Generate() string {
	return fmt.Sprintf("%06d", o.randSource.Intn(999999))
}
