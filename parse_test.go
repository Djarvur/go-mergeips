package mergeips_test

import (
	"net"
	"testing"

	"github.com/Djarvur/go-mergeips"
	"github.com/go-test/deep"

	"github.com/Djarvur/go-mergeips/internal/subnet"
)

type testParseRow struct {
	in       []string
	expected []*net.IPNet
}

var testParseData = []testParseRow{
	{
		in: []string{
			"192.168.0.3/32",
			"192.168.0.0/30",
			"192.168.0.4",
			"192.168.0.5-192.168.0.8",
		},
		expected: []*net.IPNet{
			subnet.MustParseCIDR("192.168.0.0/29", true).IPNet(),
			subnet.MustParseCIDR("192.168.0.8/32", true).IPNet(),
		},
	},
}

func TestParse(t *testing.T) {
	for _, row := range testParseData {
		nets, err := mergeips.Scan(&stringSliceScanner{data: row.in, next: -1})
		if err != nil {
			t.Error(err)
		}

		merged := mergeips.Merge(nets)
		if diff := deep.Equal(merged, row.expected); diff != nil {
			t.Errorf("got %v, expected %v: %v", merged, row.expected, diff)
		}
	}
}

//////////////

type stringSliceScanner struct {
	data []string
	next int
}

func (s *stringSliceScanner) Scan() bool {
	s.next++
	return s.next < len(s.data)
}

func (s *stringSliceScanner) Text() string {
	return s.data[s.next]
}

func (s *stringSliceScanner) Err() error {
	return nil
}

func parseCIDR(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}

	return n
}
