package util

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gosimple/slug"
)

type SlugExistsFunc func(ctx context.Context, candidate string) (bool, error)

func GenerateUniqueSlug(ctx context.Context, name string, exists SlugExistsFunc) (string, error) {
	base := slug.Make(name)
	candidate := base
	for attempt := 0; attempt < 5; attempt++ {
		exists, err := exists(ctx, candidate)
		if err != nil {
			return "", fmt.Errorf("failed to check if slug exists: %v", err)
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%s", base, randomHex(3))
	}
	return candidate, nil
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
