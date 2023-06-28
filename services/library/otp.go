package library

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
	//rand.Seed(time.Now().UnixNano())
	return &otpGenerator{randSource: randSource}
}

func (o *otpGenerator) Generate() string {
	return fmt.Sprintf("%06d", o.randSource.Intn(999999))
}
