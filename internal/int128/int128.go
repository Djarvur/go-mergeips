// Package int128 comment should be of this form
package int128

import (
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"math/bits"
	"net"
)

// Uint128 exported type should have comment or be unexported
type Uint128 struct {
	high uint64
	low  uint64
}

// Uint128FromUint64s exported func should have comment or be unexported
func Uint128FromUint64s(high, low uint64) Uint128 {
	return Uint128{high: high, low: low}
}

// ErrInvalidData exported var should have comment or be unexported
var ErrInvalidData = errors.New("invalid data")

// Uint128FromIP exported func should have comment or be unexported
func Uint128FromIP(ip net.IP) Uint128 {
	if v4 := ip.To4(); v4 != nil {
		return Uint128{
			low: uint64(binary.BigEndian.Uint32([]byte(v4))),
		}
	}

	return Uint128{
		high: binary.BigEndian.Uint64(ip[:8]),
		low:  binary.BigEndian.Uint64(ip[8:]),
	}
}

// Cmp exported func should have comment or be unexported
func (i Uint128) Cmp(j Uint128) int {
	switch {
	case i.high < j.high:
		return -1
	case i.high > j.high:
		return 1
	case i.low < j.low:
		return -1
	case i.low > j.low:
		return 1
	}

	return 0
}

// Not exported func should have comment or be unexported
func (i Uint128) Not(j Uint128) Uint128 {
	return Uint128{
		high: ^i.high,
		low:  ^i.low,
	}
}

// And exported func should have comment or be unexported
func (i Uint128) And(j Uint128) Uint128 {
	return Uint128{
		high: i.high & j.high,
		low:  i.low & j.low,
	}
}

// LeftShift exported func should have comment or be unexported
func (i Uint128) LeftShift() Uint128 {
	j := Uint128{low: i.low << 1}

	j.high = i.high << 1
	if i.low&0x8000000000000000 > 0 {
		j.high |= 1
	}

	return j
}

// NextRangeBegin exported func should have comment or be unexported
func (i Uint128) NextRangeBegin(mask Uint128) Uint128 {
	j := Uint128{high: i.high | ^mask.high, low: i.low | ^mask.low}

	if j.low < math.MaxUint64 {
		j.low++
		return j
	}

	j.low = 0
	j.high++

	return j
}

// RangeEnd exported func should have comment or be unexported
func (i Uint128) RangeEnd(mask Uint128) Uint128 {
	return Uint128{high: i.high | ^mask.high, low: i.low | ^mask.low}
}

// Jump exported func should have comment or be unexported
func (i Uint128) Jump(mask Uint128) Uint128 {
	j := Uint128{
		high: i.high | ^mask.high,
		low:  i.low | ^mask.low,
	}

	if j.low < math.MaxUint64 {
		j.low++
		return j
	}

	j.low = 0
	j.high++

	return j
}

// BigInt exported func should have comment or be unexported
func (i Uint128) BigInt() *big.Int {
	b := make([]byte, 16)

	binary.BigEndian.PutUint64(b[:8], i.high)
	binary.BigEndian.PutUint64(b[8:], i.low)

	return big.NewInt(0).SetBytes(b)
}

// IP exported func should have comment or be unexported
func (i Uint128) IP(bits int) net.IP {
	if bits == 32 {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(i.low))

		return net.IP(b)
	}

	b := make([]byte, 16)

	binary.BigEndian.PutUint64(b[:8], i.high)
	binary.BigEndian.PutUint64(b[8:], i.low)

	return net.IP(b)
}

// Ones exported func should have comment or be unexported
func (i Uint128) Ones(max int) (z int) {
	if i.low > 0 {
		z = bits.TrailingZeros64(i.low)
	} else {
		z = bits.TrailingZeros64(i.high) + 64
	}

	if z >= max {
		return 0
	}

	return max - z
}

// String exported func should have comment or be unexported
func (i Uint128) String() string {
	return i.BigInt().String()
}
