package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	pathKey := CASPathTranformsFunc(key)
	fmt.Println(pathKey)
	expectedOriginalKey := "cf5d4b01c4d9438c22c56c832f83bd3e8c6304f9"
	expectedPathName := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s  want %s", pathKey.PathName, expectedPathName)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTranformsFunc: CASPathTranformsFunc,
	}

	s := NewStore(opts)

	key := "myspecialpicture"

	data := []byte("new pc game spider man ")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTranformsFunc: CASPathTranformsFunc,
	}

	s := NewStore(opts)

	key := "myspecialpicture"

	data := []byte("new pc game spider man ")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	// b, _ := ioutil.ReadAll(r)
	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

	s.Delete(key)
}
