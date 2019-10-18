// Package masks comment should be of this form
package masks

import (
	"net"

	"github.com/Djarvur/go-mergeips/internal/bigint"
	"github.com/Djarvur/go-mergeips/internal/int128"
)

// Mask exported type should have comment or be unexported
type Mask struct {
	Mask int128.Uint128
	Size bigint.Int
}

// exported var should have comment on this block or be unexported
var (
	masksV4 = make([]Mask, 33)  // nolint: gochecknoglobals
	masksV6 = make([]Mask, 129) // nolint: gochecknoglobals
)

// Get exported func should have comment or be unexported
func Get(ones, bits int) Mask {
	switch bits {
	case 32:
		return masksV4[ones]
	case 128:
		return masksV6[ones]
	default:
		panic("invalid mask length")
	}
}

func init() { //nolint: gochecknoinits
	for ri := range masksV4 {
		masksV4[ri] = Mask{
			Mask: int128.Uint128FromIP(net.IP(net.CIDRMask(ri, 32))),
			Size: bigint.IntByBits(32).SetBit(32 - ri),
		}
	}

	for ri := range masksV6 {
		masksV6[ri] = Mask{
			Mask: int128.Uint128FromIP(net.IP(net.CIDRMask(ri, 128))),
			Size: bigint.IntByBits(128).SetBit(128 - ri),
		}
	}
}
