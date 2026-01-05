package collections

// Trie is a byte-based prefix tree used for fast path prefix checks.
type Trie struct {
	children map[byte]*Trie
	terminal bool
	wildcard bool
}

func NewTrie() *Trie {
	return &Trie{children: make(map[byte]*Trie)}
}

// Insert stores a path in the trie. When wildcard is true, any suffix will match.
func (t *Trie) Insert(path string, wildcard bool) {
	node := t
	for i := 0; i < len(path); i++ {
		ch := path[i]
		if node.children[ch] == nil {
			node.children[ch] = &Trie{children: make(map[byte]*Trie)}
		}
		node = node.children[ch]
	}
	if wildcard {
		node.wildcard = true
		return
	}
	node.terminal = true
}

// Match returns true when the path matches an inserted path or wildcard prefix.
func (t *Trie) Match(path string) bool {
	node := t
	for i := 0; i < len(path); i++ {
		if node.wildcard {
			return true
		}
		child := node.children[path[i]]
		if child == nil {
			return node.wildcard
		}
		node = child
	}
	return node.wildcard || node.terminal
}
