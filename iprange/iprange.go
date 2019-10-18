// Package iprange is to havdle ipranges, like 1.1.1.1-2.2.2.2
// IPv6 supported, of course
package iprange

import (
	"errors"
	"fmt"
	"math"
	"net"

	"github.com/Djarvur/go-mergeips/internal/int128"
	"github.com/Djarvur/go-mergeips/internal/subnet"
)

// Errors
var (
	ErrIncorrectRange = errors.New("incorrect range")
)

var (
	closedMask = int128.Uint128FromUint64s(math.MaxUint64, math.MaxUint64) // nolint: gochecknoglobals
)

// Merge returns a range as a list of subnets, as compact as possible
func Merge(begin net.IP, end net.IP) []*net.IPNet {
	bits := 128
	if begin.To4() != nil {
		bits = 32
	}

	res128 := mergeRange128(int128.Uint128FromIP(begin), int128.Uint128FromIP(end), bits)
	res := make([]*net.IPNet, 0, len(res128))

	for _, n := range res128 {
		res = append(res, n.IPNet())
	}

	return res
}

func mergeRange128(begin int128.Uint128, end int128.Uint128, bits int) (res []subnet.Subnet) {
	switch begin.Cmp(end) {
	case 1:
		panic(fmt.Errorf("%s-%s: %w", begin.IP(bits).String(), end.IP(bits).String(), ErrIncorrectRange))
	case 0:
		return []subnet.Subnet{{IP: begin, Bits: bits, Ones: bits}}
	}

	var (
		current = begin
		mask    = closedMask
	)

	for current.Cmp(end) <= 0 {
		biggerMask := mask.LeftShift()

		if current.Cmp(current.And(biggerMask)) != 0 {
			res = append(res, subnet.Subnet{IP: current, Bits: bits, Ones: mask.Ones(bits)})
			current = current.Jump(mask)
			mask = closedMask

			continue
		}

		biggerEnd := current.RangeEnd(biggerMask)

		switch biggerEnd.Cmp(end) {
		case -1:
			mask = biggerMask
			continue
		case 1:
			res = append(res, subnet.Subnet{IP: current, Bits: bits, Ones: mask.Ones(bits)})
			current = current.Jump(mask)
			mask = closedMask

			continue
		}

		res = append(res, subnet.Subnet{IP: current, Bits: bits, Ones: biggerMask.Ones(bits)})

		return res
	}

	return res
}
