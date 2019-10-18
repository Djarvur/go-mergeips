package ipnet_test

import (
	"net"
	"testing"

	"github.com/Djarvur/go-mergeips/ipnet"
	"github.com/go-test/deep"
)

type testSortRow struct {
	in       []*net.IPNet
	expected []*net.IPNet
}

var testSortData = []testSortRow{
	{
		in: []*net.IPNet{
			parseCIDR("192.168.0.4/32"),
			parseCIDR("192.168.0.0/30"),
			parseCIDR("192.168.0.8/30"),
			parseCIDR("192.168.0.0/31"),
		},
		expected: []*net.IPNet{
			parseCIDR("192.168.0.0/30"),
			parseCIDR("192.168.0.0/31"),
			parseCIDR("192.168.0.4/32"),
			parseCIDR("192.168.0.8/30"),
		},
	},
	{
		in: []*net.IPNet{
			parseCIDR("192.168.0.3/32"),
			parseCIDR("192.168.0.0/30"),
			parseCIDR("192.168.0.4/32"),
			parseCIDR("192.168.0.5/32"),
			parseCIDR("192.168.0.6/31"),
			parseCIDR("192.168.0.8/32"),
		},
		expected: []*net.IPNet{
			parseCIDR("192.168.0.0/30"),
			parseCIDR("192.168.0.3/32"),
			parseCIDR("192.168.0.4/32"),
			parseCIDR("192.168.0.5/32"),
			parseCIDR("192.168.0.6/31"),
			parseCIDR("192.168.0.8/32"),
		},
	},
}

func TestSort(t *testing.T) {
	for _, row := range testSortData {
		sorted := ipnet.Sort(row.in)
		if diff := deep.Equal(sorted, row.expected); diff != nil {
			t.Errorf("got %v, expected %v: %v", sorted, row.expected, diff)
		}
	}
}

func parseCIDR(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}

	return n
}
