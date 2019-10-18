// Package bigint comment should be of this form
package bigint

// Small exported type should have comment or be unexported
type Small int64

// SetBit exported func should have comment or be unexported
func (x Small) SetBit(i int) Int {
	return x | (1 << i)
}

// Sub exported func should have comment or be unexported
func (x Small) Sub(n Int) Int {
	return x - n.(Small)
}

// IsZero exported func should have comment or be unexported
func (x Small) IsZero() bool {
	if x < 0 {
		panic("unreachable reached")
	}

	return x == 0
}
