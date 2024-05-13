package main

var defaultStoreRoot = "defaulroot"

type PathtransformFunc func(string) string

type PathKey struct {
	Pathname string
	Filename string
}

type StoreOpts struct {
	Root              string
	PathtransformFunc PathtransformFunc
}

type Store struct{
	StoreOpts
}

func NewStore(opts StoreOpts)*Store{

	if len(opts.Root)==0{
		opts.Root=defaultStoreRoot
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string,r io.Reader)error{
	// pathKey := 
}
