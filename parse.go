// Package mergeips provides a way to convert list of IP randes definitions,
// like individual IPs, CIDR subners and begin-end ranges to the minimal list of net IPNet
package mergeips

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/Djarvur/go-mergeips/iprange"
)

// Errors
var (
	ErrInputInvalid = errors.New("invalid input")
)

// Scanner is a simple interface to support Scan() function.
// Intentionnaly compatible with bufio.Scanner
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

// Scan is used to parse source to the list of net.IPNet
func Scan(s Scanner) (res []*net.IPNet, err error) {
	for s.Scan() {
		subnets, err := Parse(s.Text(), false) // nolint: govet
		if err != nil {
			return nil, err
		}

		res = append(res, subnets...)
	}

	if err = s.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Parse parses a string to net.IPNet
// String might be in 3 forms:
// ip address itself, in v4 or v6 notation
// CIDR subnet address, v4 or v6
// IP adresses range, v4 or v6, in form begin-end
// If string is false CIDR form subnet could be defined with not-a-first addrsss in the subnet.
// Otherwise the error will be returned
func Parse(s string, strict bool) ([]*net.IPNet, error) {
	fields := strings.Split(s, "/")

	if len(fields) > 2 {
		return nil, fmt.Errorf("%q: %w", s, ErrInputInvalid)
	}

	if len(fields) == 2 {
		return parseCIDR(s, strict)
	}

	fields = strings.Split(s, "-")

	if len(fields) > 2 {
		return nil, fmt.Errorf("%q: %w", s, ErrInputInvalid)
	}

	if len(fields) == 2 {
		return parseRange(fields[0], fields[1])
	}

	return parseIP(s)
}

func parseCIDR(s string, strict bool) ([]*net.IPNet, error) {
	ip, n, err := net.ParseCIDR(s)
	if err != nil || (strict && !ip.Equal(n.IP)) {
		return nil, fmt.Errorf("%q: %w", s, ErrInputInvalid)
	}

	return []*net.IPNet{n}, nil
}

func parseRange(beginString string, endString string) ([]*net.IPNet, error) {
	var (
		begin = net.ParseIP(beginString)
		end   = net.ParseIP(endString)
	)

	if begin == nil || end == nil || (begin.To4() == nil) != (end.To4() == nil) || bytes.Compare(begin, end) > 0 {
		return nil, fmt.Errorf("%q-%q: %w", beginString, endString, ErrInputInvalid)
	}

	return iprange.Merge(begin, end), nil
}

func parseIP(s string) ([]*net.IPNet, error) {
	begin := net.ParseIP(s)
	if begin == nil {
		return nil, fmt.Errorf("%q: %w", s, ErrInputInvalid)
	}

	bits := 128
	if ipV4 := begin.To4(); ipV4 != nil {
		bits = 32
		begin = ipV4
	}

	return []*net.IPNet{{IP: begin, Mask: net.CIDRMask(bits, bits)}}, nil
}
