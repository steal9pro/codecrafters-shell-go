package autocompletition

import (
	"fmt"
	"slices"
)

type TrieNode struct {
	isEndOfWord bool
	children    map[rune]*TrieNode
}

func InitTrieNode() *TrieNode {
	return &TrieNode{
		isEndOfWord: false,
		children:    make(map[rune]*TrieNode),
	}
}

func (node *TrieNode) LoadWordsToTrie(words []string) {
	for _, w := range words {
		node.Insert(w)
	}
}

func (node *TrieNode) Insert(word string) {
	current := node
	for _, char := range word {
		if current.children[char] == nil {
			current.children[char] = InitTrieNode()
		}
		current = current.children[char]
	}
	current.isEndOfWord = true
}

func (node *TrieNode) Search(word string) bool {
	current := node
	for _, char := range word {
		if current.children[char] == nil {
			return false
		}
		current = current.children[char]
	}
	return current.isEndOfWord
}

func (node *TrieNode) StartsWith(prefix string) bool {
	current := node
	for _, char := range prefix {
		if current.children[char] == nil {
			return false
		}
		current = current.children[char]
	}
	return true
}

func (node *TrieNode) Delete(word string) bool {
	if len(word) == 0 {
		return false
	}

	current := node
	for _, char := range word {
		if current.children[char] == nil {
			return false
		}
		current = current.children[char]
	}

	if !current.isEndOfWord {
		return false
	}

	current.isEndOfWord = false
	return true
}

func (node *TrieNode) Display(idx int) {
	if node.isEndOfWord {
		fmt.Println("End of word")
	}

	for char, child := range node.children {
		fmt.Printf("node index %+v \n", idx)
		fmt.Printf("%c -> ", char)
		if child != nil {
			child.Display(idx + 1)
		}
	}
}

func (node *TrieNode) GetAllWords(prefix string) ([]string, string) {
	var words []string
	var longestCommonWord string
	var multipleCompletions bool

	current := node

	// Navigate to the prefix
	for _, char := range prefix {
		if current.children[char] == nil {
			return words, ""
		}
		current = current.children[char]
	}

	// Collect all words starting with the prefix
	node.collectWords(current, prefix, &words, &longestCommonWord, &multipleCompletions)
	slices.Sort(words)
	return words, longestCommonWord
}

func (node *TrieNode) collectWords(current *TrieNode, currentWord string, words *[]string, longestCommonWord *string, multipleCompletions *bool) {
	if len(current.children) > 1 {
		*multipleCompletions = true
	}

	if current.isEndOfWord {
		if *longestCommonWord == "" && *multipleCompletions == false {
			*longestCommonWord = currentWord
		}
		*words = append(*words, currentWord)
	}

	for char, child := range current.children {
		if child != nil {
			node.collectWords(child, currentWord+string(char), words, longestCommonWord, multipleCompletions)
		}
	}
}
