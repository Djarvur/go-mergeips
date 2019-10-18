// Package subnet_test comment should be of this form
package subnet_test

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/Djarvur/go-mergeips/internal/subnet"
)

type testRow struct {
	in       []subnet.Subnet
	expected []subnet.Subnet
}

var testData = []testRow{
	{
		in: []subnet.Subnet{
			subnet.MustParseCIDR("1.2.3.0/32", true),
			subnet.MustParseCIDR("1.2.3.1/32", true),
			subnet.MustParseCIDR("1.2.3.2/31", true),
			subnet.MustParseCIDR("1.2.3.4/30", true),
			subnet.MustParseCIDR("1.2.3.8/29", true),
			subnet.MustParseCIDR("1.2.3.16/28", true),
			subnet.MustParseCIDR("6.6.6.0/28", true),
			subnet.MustParseCIDR("6.6.6.16/29", true),
			subnet.MustParseCIDR("6.6.6.24/30", true),
			subnet.MustParseCIDR("6.6.6.28/31", true),
			subnet.MustParseCIDR("6.6.6.30/32", true),
			subnet.MustParseCIDR("6.6.6.31/32", true),
			subnet.MustParseCIDR("6.6.7.0/28", true),
			subnet.MustParseCIDR("6.6.7.16/29", true),
			subnet.MustParseCIDR("6.6.7.24/30", true),
			subnet.MustParseCIDR("6.6.7.28/31", true),
			subnet.MustParseCIDR("6.6.7.30/32", true),
			subnet.MustParseCIDR("6.6.7.31/32", true),
			subnet.MustParseCIDR("6.6.7.32/32", true),
			subnet.MustParseCIDR("6.6.7.33/32", true),
			subnet.MustParseCIDR("6.6.7.34/31", true),
			subnet.MustParseCIDR("6.6.7.36/30", true),
			subnet.MustParseCIDR("6.6.7.40/29", true),
			subnet.MustParseCIDR("6.6.7.48/28", true),
		},
		expected: []subnet.Subnet{
			subnet.MustParseCIDR("1.2.3.0/27", true),
			subnet.MustParseCIDR("6.6.6.0/27", true),
			subnet.MustParseCIDR("6.6.7.0/26", true),
		},
	},
}

// TestMerge exported func should have comment or be unexported
func TestMerge(t *testing.T) {
	for _, row := range testData {
		merged := subnet.Merge(row.in)
		if diff := deep.Equal(merged, row.expected); diff != nil {
			t.Errorf("got %v, expected %v: %v", merged, row.expected, diff)
		}
	}
}
