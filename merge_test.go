package mergeips_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Djarvur/go-mergeips"
	"github.com/Djarvur/go-mergeips/ipnet"
	"github.com/go-test/deep"
)

const testMergeDataPath = "./testdata"

var testFileName = regexp.MustCompile(`(.*/)?([^/]+)\.in\.gz$`)

type testMergeRow struct {
	path     string
	name     string
	in       []*net.IPNet
	expected []*net.IPNet
}

var testMergeData = loadFiles(testMergeDataPath, testFileName)

func TestMergeByRepeat(t *testing.T) {
	data := testMergeDataCopy(testMergeData)
	for _, row := range data {
		merged := ipnet.MergeByRepeat(ipnet.DedupSorted(ipnet.Sort(row.in)))
		if diff := deep.Equal(merged, row.expected); diff != nil {
			t.Errorf("%s%s: got %v, expected %v: %v", row.path, row.name, merged, row.expected, diff)
		}
	}
}

func TestMerge(t *testing.T) {
	data := testMergeDataCopy(testMergeData)
	for _, row := range data {
		merged := mergeips.Merge(row.in)
		if diff := deep.Equal(merged, row.expected); diff != nil {
			t.Errorf("%s%s: got %v, expected %v: %v", row.path, row.name, merged, row.expected, diff)
		}
	}
}

func TestCompare(t *testing.T) {
	data1 := testMergeDataCopy(testMergeData)
	data2 := testMergeDataCopy(testMergeData)
	for ri := range data1 {
		merged1 := mergeips.Merge(data1[ri].in)
		merged2 := ipnet.MergeByRepeat(data2[ri].in)
		if diff := deep.Equal(merged1, merged2); diff != nil {
			t.Errorf("%s%s:  %v != %v: %v", data1[ri].path, data1[ri].name, merged1, merged2, diff)
		}
	}
}

func BenchmarkMergeByRepeat(b *testing.B) {
	data := testMergeDataCopy(testMergeData)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, row := range data {
			ipnet.MergeByRepeat(row.in)
		}
	}
}

func BenchmarkMerge(b *testing.B) {
	data := testMergeDataCopy(testMergeData)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, row := range data {
			mergeips.Merge(row.in)
		}
	}
}

//////////////////////////////////

func testMergeDataCopy(in []testMergeRow) []testMergeRow {
	out := make([]testMergeRow, 0, len(in))

	for _, row := range in {
		newRow := testMergeRow{
			path:     row.path,
			name:     row.name,
			in:       make([]*net.IPNet, 0, len(row.in)),
			expected: make([]*net.IPNet, 0, len(row.expected)),
		}
		for _, n := range row.in {
			newRow.in = append(newRow.in, &net.IPNet{IP: n.IP, Mask: n.Mask})
		}
		for _, n := range row.expected {
			newRow.expected = append(newRow.expected, &net.IPNet{IP: n.IP, Mask: n.Mask})
		}
		out = append(out, newRow)
	}

	return out
}

func loadFiles(path string, name *regexp.Regexp) []testMergeRow {
	files := listFiles(path, name)
	data := make([]testMergeRow, 0, len(files))
	for _, f := range files {
		out := scanFile(f.path, f.name, ".out.gz")
		fmt.Printf("out1: %v\n", out)
		out = ipnet.MergeByRepeat(out)
		fmt.Printf("out2: %v\n", out)
		data = append(
			data,
			testMergeRow{
				path:     f.path,
				name:     f.name,
				in:       scanFile(f.path, f.name, ".in.gz"),
				expected: ipnet.MergeByRepeat(scanFile(f.path, f.name, ".out.gz")),
			},
		)
	}
	return data
}

func scanFile(p string, n string, s string) []*net.IPNet {
	f, err := os.Open(p + n + s)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}

	res, err := mergeips.Scan(bufio.NewScanner(gzr))
	if err != nil {
		panic(err)
	}

	return res
}

type fileSpec struct {
	name string
	path string
}

func listFiles(p string, r *regexp.Regexp) (res []fileSpec) {
	err := filepath.Walk(
		p,
		func(fName string, fInfo os.FileInfo, errInner error) error {
			if fInfo.Mode()&os.ModeType != 0 {
				return nil
			}
			if fields := r.FindStringSubmatch(fName); fields != nil {
				res = append(res, fileSpec{path: fields[1], name: fields[2]})
			}

			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	return res
}
