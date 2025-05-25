package autocompletition

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrieWithSpecialCharacters(t *testing.T) {
	trie := InitTrieNode()

	// Test words with special characters
	words := []string{
		"test-fuzzy",
		"hello_world",
		"file.txt",
		"my-app",
		"user@domain.com",
		"version-1.2.3",
		"config/file",
		"test",
		"fuzzy",
	}

	// Insert all words
	for _, word := range words {
		trie.Insert(word)
	}

	// Test search functionality
	for _, word := range words {
		if !trie.Search(word) {
			t.Errorf("Expected to find word '%s' in trie", word)
		}
	}

	// Test words that shouldn't exist
	nonExistentWords := []string{
		"test-fuzz",
		"hello_worl",
		"file.tx",
		"nonexistent",
	}

	for _, word := range nonExistentWords {
		if trie.Search(word) {
			t.Errorf("Did not expect to find word '%s' in trie", word)
		}
	}

	// Test prefix matching
	prefixTests := []struct {
		prefix   string
		expected bool
	}{
		{"test", true},
		{"test-", true},
		{"hello_", true},
		{"file.", true},
		{"my-", true},
		{"user@", true},
		{"version-1", true},
		{"config/", true},
		{"nonexistent", false},
		{"test-fuzzy-extra", false},
	}

	for _, test := range prefixTests {
		result := trie.StartsWith(test.prefix)
		if result != test.expected {
			t.Errorf("StartsWith('%s') = %v, expected %v", test.prefix, result, test.expected)
		}
	}
}

func TestTrieGetAllWords(t *testing.T) {
	trie := InitTrieNode()

	words := []string{
		"test",
		"test-fuzzy",
		"test-case",
		"testing",
		"hello",
		"help",
	}

	trie.LoadWordsToTrie(words)

	// Test getting all words with prefix "test"
	result := trie.GetAllWords("test")
	expected := []string{"test", "test-fuzzy", "test-case", "testing"}

	sort.Strings(result)
	sort.Strings(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetAllWords('test') = %v, expected %v", result, expected)
	}

	// Test getting all words with prefix "hel"
	result = trie.GetAllWords("hel")
	expected = []string{"hello", "help"}

	sort.Strings(result)
	sort.Strings(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetAllWords('hel') = %v, expected %v", result, expected)
	}

	// Test with non-existent prefix
	result = trie.GetAllWords("xyz")
	if len(result) != 0 {
		t.Errorf("GetAllWords('xyz') = %v, expected empty slice", result)
	}
}

func TestTrieDelete(t *testing.T) {
	trie := InitTrieNode()

	words := []string{
		"test-fuzzy",
		"test-case",
		"test",
	}

	trie.LoadWordsToTrie(words)

	// Verify words exist before deletion
	for _, word := range words {
		if !trie.Search(word) {
			t.Errorf("Word '%s' should exist before deletion", word)
		}
	}

	// Delete "test-fuzzy"
	if !trie.Delete("test-fuzzy") {
		t.Error("Failed to delete 'test-fuzzy'")
	}

	// Verify "test-fuzzy" is deleted but others remain
	if trie.Search("test-fuzzy") {
		t.Error("'test-fuzzy' should be deleted")
	}

	if !trie.Search("test-case") {
		t.Error("'test-case' should still exist")
	}

	if !trie.Search("test") {
		t.Error("'test' should still exist")
	}

	// Try to delete non-existent word
	if trie.Delete("nonexistent") {
		t.Error("Should not be able to delete non-existent word")
	}
}

func BenchmarkTrieInsert(b *testing.B) {
	trie := InitTrieNode()
	words := []string{
		"test-fuzzy",
		"hello_world",
		"file.txt",
		"my-app",
		"user@domain.com",
		"version-1.2.3",
		"config/file",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range words {
			trie.Insert(word)
		}
	}
}

func BenchmarkTrieSearch(b *testing.B) {
	trie := InitTrieNode()
	words := []string{
		"test-fuzzy",
		"hello_world",
		"file.txt",
		"my-app",
		"user@domain.com",
		"version-1.2.3",
		"config/file",
	}

	trie.LoadWordsToTrie(words)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range words {
			trie.Search(word)
		}
	}
}
