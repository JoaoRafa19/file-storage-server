package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootName = "stored_files"

// Content adressable
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:]) // [:] convert fixed array to slice

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := range sliceLen {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, fmt.Sprintf("%c", os.PathSeparator)),
		Original: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

func WithPathTransformFunc(transform PathTransformFunc) func(*Store) {
	return func(storage *Store) {
		storage.PathTransformFun = transform
	}
}

func WithRootPath(root string) func(*Store) {
	return func(s *Store) {
		s.Root = root
	}
}

type PathKey struct {
	Pathname string
	Original string
}

func (p *PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, fmt.Sprintf("%c", os.PathSeparator))
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s%c%s", p.Pathname, os.PathSeparator, p.Original)
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Original: key,
	}
}

type StoreOpts func(*Store)
type Store struct {
	// root is the folder name of the root path, containing all the filders/files of the system.
	Root             string
	PathTransformFun PathTransformFunc
}

func NewStore(opts ...StoreOpts) *Store {
	store := &Store{}

	for _, opt := range opts {
		opt(store)
	}

	if store.PathTransformFun == nil {
		store.PathTransformFun = DefaultPathTransformFunc
	}

	if len(store.Root) == 0 {
		store.Root = defaultRootName
	}

	return store
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFun(key)

	fullPathWithroot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	_, err := os.Stat(fullPathWithroot)
	return !errors.Is(err, fs.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFun(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Original)
	}()

	return os.RemoveAll(s.Root + "/" + pathKey.FirstPathName())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathkey := s.PathTransformFun(key)
	pathKeyWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathkey.FullPath())

	return os.Open(pathKeyWithRoot)
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFun(key)
	pathNameWithroot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithroot, os.ModePerm); err != nil {
		return err
	}
	fullPathWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("writen %d bytes to disk: %s", n, fullPathWithRoot)

	return nil
}
