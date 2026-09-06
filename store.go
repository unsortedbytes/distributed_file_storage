package main

import (
	// "crypto/md5"

	// "crypto/md5"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

// classic content-addressable storage(CAS) path builder
func CASPathTranformsFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) // [20] bytes -> []bytes -> [:]
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize

	// paths := make([][]bytes)
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}

	// return strings.Join(paths, "/")

}

type PathTransformsFunc func(string) PathKey

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
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

type StoreOpts struct {
	PathTranformsFunc PathTransformsFunc
}

var DefaultPathTransformsFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	PathKey := s.PathTranformsFunc(key)

	_, err := os.Stat(PathKey.FullPath())
	if err == fs.ErrNotExist {
		return false
	}

	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTranformsFunc(key)

	// if err:=os.RemoveAll(pathKey.FullPath());err != nil{
	// 	return err
	// }

	defer func() {
		log.Printf("deleted  [%s] from disk", pathKey.Filename)
	}()

	// if err := os.RemoveAll(pathKey.FullPath()); err != nil{
	// 	return err
	// }
	// return os.RemoveAll(pathKey.FullPath())
	return os.RemoveAll(pathKey.FirstPathName())
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
	pathKey := s.PathTranformsFunc(key)
	return os.Open(pathKey.FullPath())

	// f, err := os.Open(pathKey.FullPath())
	// if err != nil{
	// 	return nil, err
	// }
	// return f, nil
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTranformsFunc(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPath()
	// buf := new(bytes.Buffer)
	// io.Copy(buf, r)

	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString(filenameBytes[:])
	// // pathAndFilename
	// // filename := "somefilename"

	// pathAndFilename := pathKey.PathName + "/" + filename
	// pathAndFilename := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", n, fullPath)

	return nil
}
