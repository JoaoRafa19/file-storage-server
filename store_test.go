package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	// "github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	t.Run("DefaultPathTransform", func(t *testing.T) {
		key := "6804429f74181a63c50c3d81d733a12f14a353ff"
		pathKey := DefaultPathTransformFunc(key)
		if pathKey.PathName != key {
			t.Error(t, "have %s wants %s", pathKey, key)
		}
		if pathKey.Filename != key {
			t.Error(t, "have %s wants %s", pathKey.Filename, key)
		}
	})
	t.Run("CASPathTransformFunc", func(t *testing.T) {
		key := "momsbestpicture"
		pathKey := CASPathTransformFunc(key)
		expectedOriginalKey := "6804429f74181a63c50c3d81d733a12f14a353ff"
		expectPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
		if pathKey.PathName != expectPathName {
			t.Error(t, "have %s wants %s", pathKey.PathName, expectPathName)
		}
		if pathKey.Filename != expectedOriginalKey {
			t.Error(t, "have %s wants %s", pathKey.Filename, expectedOriginalKey)
		}
	})
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d",i)
		data := []byte("some jpeg bytes test")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("Expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, err := io.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("wants %+v, have: %+v", data, b)
		}

		if err != nil {
			t.Error(err)
		}

		// t.Run("Fails", func(t *testing.T) {
		// 	opts := StoreOpts{
		// 		PathTransformFunc: CASPathTransformFunc,
		// 	}
		// 	s := NewStore(opts)
		// 	r, err := s.Read("")
		// 	assert.Nil(t, r)

		// 	assert.NotNil(t, err, "no sutch file or directory")
		// })

		if err := s.Delete(key); err != nil {
			t.Errorf("expect not have key %s", key)
		}
	}
}

func TestDeleteKey(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)
	key := "momspecialspics"
	data := []byte("some jpeg bytes test")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
		Root:              "testStorage",
	}
	s := NewStore(opts)
	return s
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
