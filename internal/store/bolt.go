package store

import (
	"bytes"
	"time"

	"github.com/MarlikAlmighty/2miners/internal/models"
	"github.com/boltdb/bolt"
)

type Wrapper struct {
	DB *bolt.DB `json:"db,omitempty"`
}

// New database and create buckets
func New() (*Wrapper, error) {

	var buckets = []string{"users", "tokens", "ttl"}

	db, err := bolt.Open("data.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	for _, k := range buckets {

		err = db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(k))
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return &Wrapper{
		DB: db,
	}, nil
}

// Write to bucket with key and value
func (r *Wrapper) Write(uid string, value []byte) error {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.Put([]byte(uid), value)
	})
	if err != nil {
		return err
	}
	return nil
}

// WriteTTL write to bucket register token
func (r *Wrapper) WriteTTL(uid string, value []byte) error {
	if err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tokens"))
		if err := b.Put([]byte(uid), value); err != nil {
			return err
		}
		ttl := tx.Bucket([]byte("ttl"))
		if err := ttl.Put([]byte(time.Now().UTC().Format(time.RFC3339Nano)), []byte(uid)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Read from bucket with key
func (r *Wrapper) Read(bucket, uid string) ([]byte, error) {
	var result []byte
	err := r.DB.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucket))
		result = bk.Get([]byte(uid))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete uid from bucket
func (r *Wrapper) Delete(bucket, key string) error {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucket))
		if err := bk.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ReadAll read all user data from bucket `users`
func (r *Wrapper) ReadAll(bucket string) (map[string]models.User, error) {
	mp := make(map[string]models.User)
	var md models.User
	if err := r.DB.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucket))
		c := bk.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := md.UnmarshalBinary(v); err != nil {
				panic(err)
			}
			mp[string(k)] = md
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return mp, nil
}

func (r *Wrapper) Sweep(keys [][]byte) {
	var err error
	if err = r.DB.Update(func(tx *bolt.Tx) (err error) {
		bk := tx.Bucket([]byte("tokens"))
		for _, key := range keys {
			if err = bk.Delete(key); err != nil {
				panic(err.Error())
			}
		}
		return err
	}); err != nil {
		panic(err.Error())
	}
}

func (r *Wrapper) GetExpired(maxAge time.Duration) ([][]byte, error) {

	var (
		keys    [][]byte
		ttlKeys [][]byte
	)

	if err := r.DB.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("ttl")).Cursor()
		maxTime := []byte(time.Now().UTC().Add(-maxAge).Format(time.RFC3339Nano))
		for k, v := c.First(); k != nil && bytes.Compare(k, maxTime) <= 0; k, v = c.Next() {
			keys = append(keys, v)
			ttlKeys = append(ttlKeys, k)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ttl"))
		for _, key := range ttlKeys {
			if err := b.Delete(key); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return keys, nil
}

// Close database
func (r *Wrapper) Close() error {
	if err := r.DB.Close(); err != nil {
		return err
	}
	return nil
}
