package helpers

import (
	"math/rand"
	"time"
)

// newRand creates a new seeded *rand.Rand instance for use in non-global contexts.
func newRand() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}

// GenerateUUID generates a random UUID for use in the profiles.
func GenerateUUID() string {
	r := newRand() // Using a locally seeded generator
	const charset = "abcdef0123456789"
	b := make([]byte, 36)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	b[8], b[13], b[18], b[23] = '-', '-', '-', '-'
	return string(b)
}
