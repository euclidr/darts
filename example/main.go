package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/euclidr/darts"
)

type DictEntry struct {
	Key   string
	Value string
}

type ByKey []*DictEntry

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

type Lexicon struct {
	entries []*DictEntry
}

func (l *Lexicon) Add(entry *DictEntry) {
	e := DictEntry{Key: entry.Key, Value: entry.Value}
	l.entries = append(l.entries, &e)
}

func (l *Lexicon) Sort() {
	sort.Sort(ByKey(l.entries))
}

func (l *Lexicon) Keys() (keys []string) {
	keys = make([]string, len(l.entries))
	for i, e := range l.entries {
		keys[i] = e.Key
	}
	return keys
}

func (l *Lexicon) At(idx int) *DictEntry {
	return l.entries[idx]
}

func main() {
	file, err := os.Open("STPhrases.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't open file:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	var lexicon Lexicon
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		kv := strings.Split(text, "\t")
		if len(kv) != 2 {
			continue
		}
		lexicon.Add(&DictEntry{Key: kv[0], Value: kv[1]})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	lexicon.Sort()
	builder := darts.DoubleArrayBuilder{}
	builder.Build(lexicon.Keys())

	idx, matched := builder.ExactMatchSearch("系一片")
	if !matched {
		fmt.Printf("error result, %s not matched\n", "系一片")
		os.Exit(1)
	}

	e := lexicon.At(idx)
	fmt.Printf("key: %s, value: %s\n", e.Key, e.Value)
}
