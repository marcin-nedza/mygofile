package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

var defaultStoreRoot = "defaulroot"

type PathtransformFunc func(string) PathKey

func CASPathtransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

type StoreOpts struct {
	Root              string
	PathtransformFunc PathtransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {

	if len(opts.Root) == 0 {
		opts.Root = defaultStoreRoot
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store)Write(key string, r io.Reader)(int64,error) {
	return s.writeStream(key,r)
}

func (s *Store) writeStream(key string, r io.Reader) (int64,error) {
	pathKey := s.PathtransformFunc(key)
	pathnamewithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)

	err := os.MkdirAll(pathnamewithRoot, os.ModePerm)
	if err != nil {
		return 0,err
	}

	fullpathwithroot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	f, err := os.Create(fullpathwithroot)
	if err != nil {
		return 0,err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return 0,err
	}
	//write to disk
	return n,nil
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}
