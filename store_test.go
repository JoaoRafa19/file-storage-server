package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newStore() *Store {
	return nil
}

func tearDown(t *testing.T, s *Store) {
	assert.Nil(t, s.Clear())
}

func TestPathTransformFunc(t *testing.T) {
	key := "yoursmamagiantpicture"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)
	expectedFilename := "711b625cc740fc466556ecf598adf8e3157be229"
	expecttedPathName := "711b6/25cc7/40fc4/66556/ecf59/8adf8/e3157/be229"
	assert.Equal(t, pathname.Pathname, expecttedPathName)
	assert.Equal(t, pathname.Original, expectedFilename)
}

func TestStoreCasPath(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)

	defer tearDown(t, s)

	data := bytes.NewReader([]byte("mypicture"))
	err := s.writeStream("mypicture", data)

	assert.Nil(t, err)

	s.Delete("mypicture")
}

func TestStore(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)

	defer tearDown(t, s)

	for i := range 50 {
		key := fmt.Sprintf("%d-key", i)
		data := []byte(key)
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expect to not have the key")
		}
	}

}

func TestStoreRead(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)
	defer tearDown(t, s)
	key := "mypicture"
	data := []byte("some jpeg bytes")
	reader := bytes.NewReader(data)
	err := s.writeStream(key, reader)

	assert.Nil(t, err)

	r, err := s.Read(key)
	assert.Nil(t, err)

	b, err := io.ReadAll(r)
	// print file
	fmt.Println(string(b))
	assert.Nil(t, err)

	// assert.Equal(t, data, b)

	s.Delete(key)

}

func TestStoreDeleteKey(t *testing.T) {
	s := NewStore(
		WithPathTransformFunc(CASPathTransformFunc),
	)

	defer tearDown(t, s)

	key := "mypicture"
	data := []byte("some jpeg bytes")
	reader := bytes.NewReader(data)
	err := s.writeStream(key, reader)

	assert.Nil(t, err)

	assert.Nil(t, s.Delete(key))
}
