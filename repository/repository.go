package repository

import (
	"github.com/boltdb/bolt"
)

type Repository struct {
	db *bolt.DB
}

func NewRepository(conn *bolt.DB) *Repository {
	err := conn.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Bucket"))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &Repository{db: conn}
}

func (r Repository) Get(key string) string {
	var result string

	r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Bucket"))

		v := b.Get([]byte(key))

		if len(v) > 0 {
			result = string(v)
		}

		return nil
	})

	return result
}

func (r *Repository) Put(key, value string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Bucket"))

		return b.Put([]byte(key), []byte(value))
	})
}

func (r *Repository) Delete(key string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Bucket"))

		return b.Delete([]byte(key))
	})
}