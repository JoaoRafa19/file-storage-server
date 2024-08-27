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

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashString) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, string(os.PathSeparator)),
		Filename: hashString,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s%c%s", p.PathName, os.PathSeparator, p.Filename)
}

const defaultFoldername = "storager"

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
	Root              string // Root Storage folder name
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Filename: key,
		PathName: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultFoldername
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := new(bytes.Buffer)

	_, err = io.Copy(buff, f)
	return buff, err
}

func (s *Store) Has(Key string) bool {
	pathKey := s.PathTransformFunc(Key)

	fullPathWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)

	return !errors.Is(err, fs.ErrNotExist)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	fullPath := s.PathTransformFunc(key)
	pathkeyWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, fullPath.FullPath())
	return os.Open(pathkeyWithRoot)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FirstPathName())
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.PathName)
	if err := os.MkdirAll(pathWithRoot, os.ModePerm); err != nil {
		return err
	}

	// buff := new(bytes.Buffer)

	fullPathWithRoot := fmt.Sprintf("%s%c%s", s.Root, os.PathSeparator, pathKey.FullPath())

	// filenameBytes := md5.Sum(buff.Bytes())
	// filename := hex.EncodeToString(filenameBytes[:])

	// fullPath := pathKey.PathName + string(os.PathSeparator) + filename

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("Writen (%d) bytes to disk: %s\n", n, fullPathWithRoot)

	return nil
}
