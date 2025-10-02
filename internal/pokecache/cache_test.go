package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
	}

	for i, cse := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(cse.key, cse.val)
			got, ok := cache.Get(cse.key)
			if !ok {
				t.Fatalf("expected to find key %q", cse.key)
			}
			if string(got) != string(cse.val) {
				t.Fatalf("expected %q, got %q", string(cse.val), string(got))
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseInterval = 5 * time.Millisecond
	const waitTime = baseInterval + 5*time.Millisecond

	cache := NewCache(baseInterval)
	cache.Add("https://example.com", []byte("testdata"))

	if _, ok := cache.Get("https://example.com"); !ok {
		t.Fatalf("expected key to exist immediately after Add")
	}

	time.Sleep(waitTime)

	if _, ok := cache.Get("https://example.com"); ok {
		t.Fatalf("expected key to be reaped (not found) after wait")
	}
}
