package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func TestEntryListSave(t *testing.T) {
	rawEntries := []*entry{
		{"b", 10}, {"a", 20}, {"c", 15},
	}
	entries := entryList(rawEntries)

	entriesFile, _ := ioutil.TempFile("", "testEntries")
	fileName := entriesFile.Name()
	defer os.Remove(fileName)

	entries.Save(fileName)

	entriesFile, err := os.Open(fileName)
	if err != nil {
		t.Error("Failed to read the saved entries file.", err)
	}
	scanner := bufio.NewScanner(entriesFile)
	results := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)
	}
	if len(results) != len(entries) {
		t.Errorf("Incorrect number of entries saved: %q", results)
	}
	for i, e := range entries {
		if results[i] != e.String() {
			t.Errorf("Entry %d saved incorrectly: %q", i, results[i])
		}
	}
}

func TestEntryListSort(t *testing.T) {
	rawEntries := []*entry{
		{"b", 10},
		{"a", 20},
		{"c", 15},
	}
	entries := entryList(rawEntries)
	entries.Sort()
	expected := []string{"a", "c", "b"}
	for i, e := range entries {
		if expected[i] != e.val {
			t.Errorf("Item %d not in place, expected %s, got %s", i, expected[i], e.val)
		}
	}
}

func TestEntryListFilter(t *testing.T) {
	entries := entryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	nonZeroScore := func(e *entry) bool { return e.score > 0 }
	nonZeroEntries := entries.Filter(nonZeroScore)
	if len(nonZeroEntries) != 1 {
		t.Errorf("Entries not filtered correctly: %v", nonZeroEntries)
	}
	if nonZeroEntries[0].val != "/path_b" {
		t.Errorf("Incorrect entry left after filtering: %v", nonZeroEntries)
	}
}

func TestEntryListUpdate(t *testing.T) {
	entries := entryList{
		&entry{"/path_b", 10},
		&entry{"/path_a", 0},
	}
	entries = entries.Update("/path_a", 1)
	if entries[0].score != 10 || entries[1].score != 1 {
		t.Errorf("Invalid update: %v", entries)
	}
	entries = entries.Update("/path_c", 1)
	if len(entries) != 3 {
		t.Errorf("New entry not created: %d", len(entries))
	}
}

func TestEntryListAge(t *testing.T) {
	entries := entryList{
		&entry{"a", 20},
		&entry{"b", 10},
		&entry{"c", 0},
	}
	entries.Age()
	expected := []float64{18.0, 9.0, 0}
	for i, e := range entries {
		if e.score != expected[i] {
			t.Errorf("Score not updated correctly, expect %f, get %f", expected[i], e.score)
		}
	}
}

func TestString(t *testing.T) {
	e := &entry{"/etc/init", 10.1234}
	if e.String() != "10.12\t/etc/init" {
		t.Errorf("Wrong string representation: %s", e.String())
	}
}

func TestUpdateEntryScore(t *testing.T) {
	e := &entry{"/etc/init", 0}
	e.updateScore(10)
	if e.score != 10 {
		t.Errorf("Entity score is wrong: %f", e.score)
	}
	e.updateScore(10)
	if e.score-14.14 < 0.001 {
		t.Errorf("Entity score is wrong: %f", e.score)
	}
}
