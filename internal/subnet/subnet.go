// Package subnet comment should be of this form
package subnet

import (
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/Djarvur/go-mergeips/internal/int128"
	"github.com/Djarvur/go-mergeips/internal/masks"
)

// Subnet exported type should have comment or be unexported
type Subnet struct {
	IP   int128.Uint128
	Ones int
	Bits int
}

// Include exported func should have comment or be unexported
func (s Subnet) Include(b Subnet) bool {
	return s.Bits == b.Bits && s.IP.Cmp(b.IP.And(masks.Get(s.Ones, s.Bits).Mask)) == 0
}

// String exported func should have comment or be unexported
func (s Subnet) String() string {
	return s.IPNet().String()
}

// FromIPNet exported func should have comment or be unexported
func FromIPNet(n *net.IPNet) Subnet {
	ones, bits := n.Mask.Size()

	return Subnet{
		IP:   int128.Uint128FromIP(n.IP),
		Ones: ones,
		Bits: bits,
	}
}

// IPNet exported func should have comment or be unexported
func (s Subnet) IPNet() *net.IPNet {
	return &net.IPNet{
		IP:   s.IP.IP(s.Bits),
		Mask: net.CIDRMask(s.Ones, s.Bits),
	}
}

// Sort exported func should have comment or be unexported
func Sort(ips []Subnet) []Subnet {
	sort.Slice(ips, func(i, j int) bool { return ips[i].Less(ips[j]) })
	return ips
}

// Less exported func should have comment or be unexported
func (s Subnet) Less(b Subnet) bool {
	if s.Bits != b.Bits {
		return s.Bits < b.Bits
	}

	if cmp := s.IP.Cmp(b.IP); cmp != 0 {
		return cmp < 0
	}

	return s.Ones < b.Ones
}

// Mask exported func should have comment or be unexported
func (s Subnet) Mask() masks.Mask {
	return masks.Get(s.Ones, s.Bits)
}

// DedupSorted exported func should have comment or be unexported
func DedupSorted(ips []Subnet) []Subnet {
	j := 0

	for i := 1; i < len(ips); i++ {
		if ips[j].Include(ips[i]) {
			continue
		}
		j++

		ips[j] = ips[i]
	}

	return ips[:j+1]
}

// MergePairs exported func should have comment or be unexported
func MergePairs(ips []Subnet) []Subnet {
	j := 0

	for i := 1; i < len(ips); i++ {
		bigger, ok := biggerSubnet(ips[j])
		if ok && ips[j].Ones == ips[i].Ones && bigger.Include(ips[i]) {
			ips[j] = bigger
			continue
		}
		j++

		ips[j] = ips[i]
	}

	return ips[:j+1]
}

// MergeSorted exported func should have comment or be unexported
func MergeSorted(ips []Subnet) []Subnet {
	for newips := MergePairs(ips); len(newips) != len(ips); newips = MergePairs(ips) {
		ips = newips
	}

	return ips
}

// Merge exported func should have comment or be unexported
func Merge(ips []Subnet) []Subnet {
	ips = DedupSorted(Sort(ips))
	for newips := MergePairs(ips); len(newips) != len(ips); newips = MergePairs(ips) {
		ips = newips
	}

	return ips
}

func biggerSubnet(s Subnet) (Subnet, bool) {
	if s.Ones == 0 {
		return s, true
	}

	bigger := Subnet{IP: s.IP, Bits: s.Bits, Ones: s.Ones - 1}

	if s.IP.Cmp(s.IP.And(masks.Get(bigger.Ones, bigger.Bits).Mask)) != 0 {
		return s, false
	}

	return bigger, true
}

// ErrIncorrectCIDR exported var should have comment or be unexported
var ErrIncorrectCIDR = errors.New("CIDR definition incorrect")

// ParseCIDR exported func should have comment or be unexported
func ParseCIDR(s string, strict bool) (Subnet, error) {
	ip, n, err := net.ParseCIDR(s)
	if err != nil {
		return Subnet{}, err
	}

	if strict && !ip.Equal(n.IP) {
		return Subnet{}, fmt.Errorf("%q expected, got %q: %w", n.String(), s, ErrIncorrectCIDR)
	}

	return FromIPNet(n), nil
}

// MustParseCIDR exported func should have comment or be unexported
func MustParseCIDR(s string, strict bool) Subnet {
	n, err := ParseCIDR(s, strict)
	if err != nil {
		panic(err)
	}

	return n
}
