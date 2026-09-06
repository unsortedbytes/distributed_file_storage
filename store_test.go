package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTranformsFunc: CASPathTranformsFunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

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

// func TestStoreDeleteKey(t *testing.T) {
// 	opts := StoreOpts{
// 		PathTranformsFunc: CASPathTranformsFunc,
// 	}

// 	s := NewStore(opts)

// 	key := "myspecialpicture"

// 	data := []byte("new pc game spider man ")
// 	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
// 		t.Error(err)
// 	}

// 	if err := s.Delete(key); err != nil {
// 		t.Error(err)
// 	}
// }

func TestStore(t *testing.T) {
	// opts := StoreOpts{
	// 	PathTranformsFunc: CASPathTranformsFunc,
	// }

	// s := NewStore(opts)

	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {

		// key := "foodbar"
		key := fmt.Sprintf("foo_%d", i)

		data := []byte("new pc game spider man ")
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

		// b, _ := ioutil.ReadAll(r)
		b, _ := io.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		fmt.Println(string(b))

		// s.Delete(key)
		// if err := s.Delete(key); err != nil {
		// 	t.Error(err)
		// }

		// if ok := s.Has(key); !ok {
		// 	t.Errorf("expected to Not have key %s", key)
		// }
	}
}
