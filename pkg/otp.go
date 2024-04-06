package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

type OtpGenerator interface {
	Generate() string
}

type otpGenerator struct {
	randSource *rand.Rand
}

func NewOTPGenerator() OtpGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	randSource := rand.New(source)
	return &otpGenerator{randSource: randSource}
}

func (o *otpGenerator) Generate() string {
	return fmt.Sprintf("%06d", o.randSource.Intn(999999))
}
