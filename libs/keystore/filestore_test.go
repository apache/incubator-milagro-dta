package keystore

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileStore(t *testing.T) {
	keys := []Key{{"key1", []byte{1}}, {"key2", []byte{1, 2}}, {"key3", []byte{1, 2, 3}}}
	keynames := []string{"key1", "key2", "key3"}

	fn := tmpFileName()
	defer func() {
		if err := os.Remove(fn); err != nil {
			t.Logf("Warning! Temp file could not be deleted (%v): %v", err, fn)
		}
	}()

	fs, err := NewFileStore(fn)
	if err != nil {
		t.Fatal(err)
	}
	fs.Set(keys...)
	skeys, err := fs.Get(keynames...)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range keys {
		if !bytes.Equal(k.Key, skeys[k.Name]) {
			t.Errorf("Key not match: %v. Expected: %v, Found: %v", k.Name, k.Key, skeys[k.Name])
		}
	}

	fs1, err := NewFileStore(fn)
	if err != nil {
		t.Fatal(err)
	}
	s1keys, err := fs1.Get(keynames...)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range keys {
		if !bytes.Equal(k.Key, s1keys[k.Name]) {
			t.Errorf("Key not match: %v. Expected: %v, Found: %v", k.Name, k.Key, s1keys[k.Name])
		}
	}
}

func tmpFileName() string {
	rnd := make([]byte, 8)
	rand.Read(rnd)
	ts := time.Now().UnixNano()

	return filepath.Join(os.TempDir(), fmt.Sprintf("keystore-%v-%x.tmp", ts, rnd))
}
