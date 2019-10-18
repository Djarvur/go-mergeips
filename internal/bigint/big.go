// Package bigint comment should be of this form
package bigint

import "math/big"

// Big exported type should have comment or be unexported
type Big struct {
	*big.Int
}

var bigInt0 = big.NewInt(0) // nolint: gochecknoglobals

// SetBit exported func should have comment or be unexported
func (x Big) SetBit(i int) Int {
	return Big{big.NewInt(0).SetBit(bigInt0, i, 1)}
}

// Sub exported func should have comment or be unexported
func (x Big) Sub(n Int) Int {
	return Big{big.NewInt(0).Sub(x.Int, n.(Big).Int)}
}

// IsZero exported func should have comment or be unexported
func (x Big) IsZero() bool {
	cmp := x.Int.Cmp(bigInt0)
	if cmp < 0 {
		panic("unreachable reached")
	}

	return cmp == 0
}
