package store

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type with struct {
	db     *bolt.DB
	bucket []byte
}

func (w *with) get(k []byte) (out []byte, err error) {
	err = w.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(w.bucket)
		out = bucket.Get(k)
		return nil
	})
	if out == nil {
		return nil, errors.New("not exist")
	}
	return out, err
}
func (w *with) set(k []byte, v []byte) error {
	err := w.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(w.bucket)
		return bucket.Put(k, v)
	})
	return err
}
func (w *with) del(k []byte) error {
	err := w.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(w.bucket)
		return bucket.Delete(k)
	})
	return err
}
func (w *with) Get(k string) (string, error) {
	ret, err := w.get([]byte(k))
	return string(ret), err
}
func (w *with) Getx(k string, v interface{}) error {
	ret, err := w.get([]byte(k))
	if err != nil {
		return err
	}
	return json.Unmarshal(ret, v)
}

func (w *with) Set(k, v string) error {
	return w.set([]byte(k), []byte(v))
}

func (w *with) Setx(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return w.set([]byte(k), b)
}

func (w *with) Del(k string) error {
	return w.del([]byte(k))
}
func (w *with) first() (out []byte, val []byte, err error) {
	err = w.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(w.bucket)
		out, val = bucket.Cursor().First()
		return nil
	})
	if out == nil && val == nil {
		err = errors.New("empty")
		return
	}
	return
}
func (w *with) First() (string, string, error) {
	key, val, err := w.first()
	return string(key), string(val), err
}

func (w *with) Firstx(v interface{}) (string, error) {
	key, val, err := w.first()
	if err != nil {
		return "", err
	}
	return string(key), json.Unmarshal(val, v)
}

// func (w *with) eachKeys() []string {
// 	w.db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket(w.bucket)
// 		c := b.Cursor()
// 		for k, v := c.First(); k != nil; k, v = c.Next() {
// 			fmt.Printf("key=%s, value=%s\n", k, v)
// 		}

// 		return nil
// 	})
// }

func (w *with) Each(fn func(k string, v []byte)) {
	w.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(w.bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fn(string(k), v)
		}
		return nil
	})
}
