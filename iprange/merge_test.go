package iprange_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/Djarvur/go-mergeips/iprange"
	"github.com/go-test/deep"
)

type testRowMergeRange struct {
	begin    net.IP
	end      net.IP
	expected []*net.IPNet
}

var testTableMergeRangeV4 = []testRowMergeRange{
	{begin: net.ParseIP("0.0.0.0"), end: net.ParseIP("255.255.255.255"), expected: []*net.IPNet{parseCIDR("0.0.0.0/0")}},
	{begin: net.ParseIP("0.0.0.0"), end: net.ParseIP("0.0.0.0"), expected: []*net.IPNet{parseCIDR("0.0.0.0/32")}},
	{begin: net.ParseIP("255.255.255.255"), end: net.ParseIP("255.255.255.255"), expected: []*net.IPNet{parseCIDR("255.255.255.255/32")}},
	{begin: net.ParseIP("192.168.0.7"), end: net.ParseIP("192.168.0.7"), expected: []*net.IPNet{parseCIDR("192.168.0.7/32")}},
	{begin: net.ParseIP("192.168.0.0"), end: net.ParseIP("192.168.0.255"), expected: []*net.IPNet{parseCIDR("192.168.0.0/24")}},
	{begin: net.ParseIP("192.168.0.7"), end: net.ParseIP("192.168.0.22"),
		expected: []*net.IPNet{
			parseCIDR("192.168.0.7/32"),
			parseCIDR("192.168.0.8/29"),
			parseCIDR("192.168.0.16/30"),
			parseCIDR("192.168.0.20/31"),
			parseCIDR("192.168.0.22/32"),
		},
	},
	{begin: net.ParseIP("0.0.0.0"), end: net.ParseIP("0.0.0.10"),
		expected: []*net.IPNet{
			parseCIDR("0.0.0.0/29"),
			parseCIDR("0.0.0.8/31"),
			parseCIDR("0.0.0.10/32"),
		},
	},
	{begin: net.ParseIP("255.255.255.247"), end: net.ParseIP("255.255.255.255"),
		expected: []*net.IPNet{
			parseCIDR("255.255.255.247/32"),
			parseCIDR("255.255.255.248/29"),
		},
	},
}

func TestMergeRange(t *testing.T) {
	for _, r := range testTableMergeRangeV4 {
		subnets := iprange.Merge(r.begin, r.end)
		if diff := deep.Equal(subnets, r.expected); diff != nil {
			t.Errorf("got %v, expected %v: %v", subnets, r.expected, diff)
		}
		fmt.Printf("subnets: %v\n", subnets)
	}
}

func BenchmarkMergeRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, r := range testTableMergeRangeV4 {
			iprange.Merge(r.begin, r.end)
		}
	}
}

func parseCIDR(s string) *net.IPNet {
	ip, n, err := net.ParseCIDR(s)
	if err != nil || !ip.Equal(n.IP) {
		return nil
	}

	return n
}
