package ipnet_test

import (
	"net"
	"testing"

	"github.com/Djarvur/go-mergeips/ipnet"
	"github.com/go-test/deep"
)

type testMergeRow struct {
	in       []*net.IPNet
	expected []*net.IPNet
}

var testMergeData = []testMergeRow{
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
			parseCIDR("192.168.0.0/29"),
			parseCIDR("192.168.0.8/32"),
		},
	},
}

func TestMerge(t *testing.T) {
	for _, row := range testMergeData {
		out := ipnet.Merge(row.in)
		if diff := deep.Equal(out, row.expected); diff != nil {
			t.Errorf("got %v, expected %v: %v", out, row.expected, diff)
		}
	}
}
