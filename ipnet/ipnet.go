// Package ipnet provides some useful methods to handle net.IPNet lists
package ipnet

import (
	"bytes"
	"fmt"
	"net"
	"sort"

	"github.com/Djarvur/go-mergeips/internal/bigint"
	"github.com/Djarvur/go-mergeips/internal/masks"
)

// MergeByRepeat is a wrapper around MergeSorted
func MergeByRepeat(nets []*net.IPNet) []*net.IPNet {
	return MergeSortedByRepeat(DedupSorted(Sort(nets)))
}

// MergePairs merges all the suitable pairs of subnets in the net.IPNet list
func MergePairs(nets []*net.IPNet) []*net.IPNet {
	j := 0

	for i := 1; i < len(nets); i++ {
		bigger, _ := biggerIPNet(nets[j])
		if bigger != nil && bigger.Contains(nets[i].IP) && bytes.Equal(nets[j].Mask, nets[i].Mask) {
			nets[j] = bigger
			continue
		}
		j++

		nets[j] = nets[i]
	}

	return nets[:j+1]
}

// MergeSortedByRepeat is repeating MergePairs as long as it does merge anything
func MergeSortedByRepeat(nets []*net.IPNet) []*net.IPNet {
	for newnets := MergePairs(nets); len(newnets) != len(nets); newnets = MergePairs(nets) {
		nets = newnets
	}

	return nets
}

// MergeSorted is merging previously sorted and de-duped list of net.IPNet to the smallest possible form
func MergeSorted(nets []*net.IPNet) []*net.IPNet {
	return mergeSortedBig(nets)
}

func mergeSortedBig(nets []*net.IPNet) []*net.IPNet {
	panic("does not work")
	if len(nets) == 0 {
		return nets
	}

	doneUpTo := 0
	bigger, toFill := biggerIPNet(nets[doneUpTo])
	fullyConsumed := 0

	for lastChecked := 1; lastChecked < len(nets); lastChecked++ {
		if bigger != nil && bigger.Contains(nets[lastChecked].IP) {
			toFill = toFill.Sub(masks.Get(nets[lastChecked].Mask.Size()).Size)
			if toFill.IsZero() {
				fmt.Printf("%v %v %v %v %v %v\n", doneUpTo, fullyConsumed, lastChecked, nets[doneUpTo], nets[lastChecked], bigger)
				nets[doneUpTo] = bigger
				bigger, toFill = biggerIPNet(nets[doneUpTo])
				fullyConsumed = lastChecked
			}

			continue
		}

		if fullyConsumed != lastChecked {
			lastChecked = fullyConsumed
		}

		doneUpTo++
		nets[doneUpTo] = nets[fullyConsumed]
		bigger, toFill = biggerIPNet(nets[doneUpTo])
	}

	if fullyConsumed < len(nets) {
		return append(nets[:doneUpTo+1], nets[fullyConsumed:]...)
	}

	return nets[:doneUpTo+1]
}

func biggerIPNet(n *net.IPNet) (*net.IPNet, bigint.Int) {
	ones, bits := n.Mask.Size()
	if ones == 0 {
		return n, masks.Get(ones, bits).Size
	}

	biggerMask := net.CIDRMask(ones-1, bits)
	if !n.IP.Equal(n.IP.Mask(biggerMask)) {
		return nil, masks.Get(bits, bits).Size
	}

	return &net.IPNet{IP: n.IP, Mask: biggerMask}, masks.Get(ones, bits).Size
}

// Sort sorts lust of net.IPNet and return it
// IPv4 goes first, bigger mask goes first
func Sort(nets []*net.IPNet) []*net.IPNet {
	sort.Slice(nets, func(i, j int) bool { return Less(nets[i], nets[j]) })
	return nets
}

// Less is comparing to net.IPNet
// To be used with Sort()
func Less(a, b *net.IPNet) bool {
	if cmp := bytes.Compare(a.IP, b.IP); cmp != 0 {
		return cmp < 0
	}

	aOnes, _ := a.Mask.Size()
	bOnes, _ := b.Mask.Size()

	return aOnes < bOnes
}

// DedupSorted removes all the identical or included-in-bigger-one-presented sublens from the list
func DedupSorted(nets []*net.IPNet) []*net.IPNet {
	if len(nets) == 0 {
		return nets
	}

	j := 0

	for i := 1; i < len(nets); i++ {
		if nets[j].Contains(nets[i].IP) {
			continue
		}
		j++

		nets[j] = nets[i]
	}

	return nets[:j+1]
}
