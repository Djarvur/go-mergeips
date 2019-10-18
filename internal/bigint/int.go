// Package bigint comment should be of this form
package bigint

import (
	"math/big"
)

// Int exported type should have comment or be unexported
type Int interface {
	SetBit(i int) Int
	Sub(n Int) Int
	IsZero() bool
}

// IntByBits exported func should have comment or be unexported
func IntByBits(bits int) Int {
	switch bits {
	case 32:
		return Small(0)
	case 128:
		return Big{big.NewInt(0)}
	default:
		panic("incottect bits")
	}
}
