package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type MFACode struct {
	Code      string
	ExpiresAt time.Time
}

func GenerateMFA() *MFACode {
	// Generate a random 6-digit code (000000 to 999999)
	max := big.NewInt(1000000) // 10^6
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to a default in case of error
		return &MFACode{
			Code:      "000000",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
	}

	// Format as 6-digit string with leading zeros
	code := fmt.Sprintf("%06d", n.Int64())

	return &MFACode{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func (m *MFACode) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

func (m *MFACode) TimeRemaining() time.Duration {
	remaining := time.Until(m.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
