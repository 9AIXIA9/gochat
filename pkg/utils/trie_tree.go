package utils

type TrieNode struct {
	children map[byte]*TrieNode
	terminal bool
	wildcard bool
}

func NewTrie() *TrieNode {
	return &TrieNode{
		children: make(map[byte]*TrieNode),
	}
}

func (n *TrieNode) Insert(path string, wildcard bool) {
	node := n
	for i := 0; i < len(path); i++ {
		ch := path[i]
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[byte]*TrieNode)}
		}
		node = node.children[ch]
	}
	if wildcard {
		node.wildcard = true
		return
	}
	node.terminal = true
}

func (n *TrieNode) Match(path string) bool {
	node := n
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
