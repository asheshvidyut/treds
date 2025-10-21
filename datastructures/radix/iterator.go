package radix

import (
	"bytes"
	"regexp"
)

// Iterator is used to iterate over a set of nodes
// in pre-order
type Iterator struct {
	node           *Node
	stack          []edges
	iterLeafNode   *LeafNode
	iterCounter    int
	key            []byte
	seekLowerBound bool
	pattern        *regexp.Regexp
	regexMatch     bool
}

func (i *Iterator) PatternMatch(rx *regexp.Regexp) {
	i.regexMatch = true
	i.pattern = rx
}

// SeekPrefixWatch is used to seek the iterator to a given prefix
// and returns the watch channel of the finest granularity
func (i *Iterator) SeekPrefixWatch(prefix []byte) (watch <-chan struct{}) {
	// Wipe the stack
	i.seekLowerBound = false
	i.stack = nil
	i.key = prefix
	n := i.node
	search := prefix
	i.iterLeafNode = i.node.minLeaf
	i.iterCounter = i.node.leavesInSubtree
	for {
		// Check for key exhaustion
		if len(search) == 0 {
			i.node = n
			i.iterLeafNode = i.node.minLeaf
			i.iterCounter = i.node.leavesInSubtree
			return
		}

		// Look for an edge
		_, n = n.getEdge(search[0])
		if n == nil {
			i.node = nil
			i.iterLeafNode = nil
			i.iterCounter = 0
			return
		}

		// Consume the search prefix
		if bytes.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]

		} else if bytes.HasPrefix(n.prefix, search) {
			i.node = n
			i.iterLeafNode = i.node.minLeaf
			i.iterCounter = i.node.leavesInSubtree
			return
		} else {
			i.node = nil
			i.iterLeafNode = nil
			i.iterCounter = 0
			return
		}
	}
}

// SeekPrefix is used to seek the iterator to a given prefix
func (i *Iterator) SeekPrefix(prefix []byte) {
	i.SeekPrefixWatch(prefix)
}

// Next returns the next node in order
func (i *Iterator) Next() ([]byte, interface{}, bool) {

	var zero interface{}

	if i.regexMatch {
		for i.iterLeafNode != nil {
			matched := i.pattern.MatchString(string(i.iterLeafNode.key))
			if i.iterLeafNode != nil && matched {
				res := i.iterLeafNode
				i.iterLeafNode = i.iterLeafNode.GetNextLeaf()
				return res.key, res.val, true
			} else {
				i.iterLeafNode = i.iterLeafNode.GetNextLeaf()
			}
		}
	} else {
		var zero interface{}

		i.iterCounter--

		if i.iterCounter < 0 {
			return nil, zero, false
		}

		key := i.iterLeafNode.key
		val := i.iterLeafNode.val
		i.iterLeafNode = i.iterLeafNode.nextLeaf
		return key, val, true
	}

	i.iterLeafNode = nil
	i.node = nil
	return nil, zero, false
}
