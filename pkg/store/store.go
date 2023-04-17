package store

import (
	"github.com/boltdb/bolt"
)

type Store struct {
	db     *bolt.DB
	bucket []byte
}

func New(dbname string) (*Store, error) {
	store := &Store{}
	store.bucket = []byte("zroot")
	db, err := bolt.Open(dbname, 0644, nil)
	if err != nil {
		return nil, err
	}
	store.db = db
	if store.init(store.bucket) != nil {
		return nil, err
	}
	return store, nil
}
func (s *Store) With(bucket string) *with {
	b := []byte(bucket)
	s.init(b)
	return s.with(b)
}
func (s *Store) with(bucket []byte) *with {
	w := &with{
		bucket: bucket,
		db:     s.db,
	}
	return w
}
func (s *Store) init(bucket []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})
	return err
}
func (s *Store) Get(k string) (string, error) {
	return s.with(s.bucket).Get(k)
}
func (s *Store) Getx(k string, v interface{}) error {
	return s.with(s.bucket).Getx(k, v)
}
func (s *Store) Set(k, v string) error {
	return s.with(s.bucket).Set(k, v)
}

func (s *Store) Setx(k string, v interface{}) error {
	return s.with(s.bucket).Setx(k, v)
}

func (s *Store) Del(k string) error {
	return s.with(s.bucket).Del(k)
}
func (s *Store) Each(fn func(k string, v []byte)) {
	s.with(s.bucket).Each(fn)
}
func (s *Store) First() (string, string, error) {
	return s.with(s.bucket).First()
}

func (s *Store) Firstx(v interface{}) (string, error) {
	return s.with(s.bucket).Firstx(v)
}

var store *Store

func init() {
	var err error
	store, err = New("/etc/.data.db")
	if err != nil {
		panic(err)
	}
}

func Get(k string) (string, error) {
	return store.Get(k)
}

func Getx(k string, v interface{}) error {
	return store.Getx(k, v)
}
func Set(k, v string) error {
	return store.Set(k, v)
}
func With(bucket string) *with {
	return store.With(bucket)
}
func Setx(k string, v interface{}) error {
	return store.Setx(k, v)
}
func Del(k string) error {
	return store.Del(k)
}

func First() (string, string, error) {
	return store.First()
}

func Firstx(v interface{}) (string, error) {
	return store.Firstx(v)
}
