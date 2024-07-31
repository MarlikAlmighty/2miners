package store

import (
	"github.com/MarlikAlmighty/2miners/internal/models"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {

	wrapper, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, wrapper.DB)

	var (
		bk  *bolt.Bucket
		md  models.User
		uid uuid.UUID
		b   []byte
	)

	// Check if buckets are created
	err = wrapper.DB.Update(func(tx *bolt.Tx) error {

		bk, err = tx.CreateBucketIfNotExists([]byte("users"))
		assert.NoError(t, err)

		uid = uuid.New()
		md.Email = "cryptocoin62@gmail.com"
		md.Pass = "b2e8d6ccb3b7df55461668f56b4abaeb"
		md.UID = uid.String()
		md.Root = true

		b, err = md.MarshalBinary()
		if err != nil {
			panic(err)
		}

		if err = bk.Put([]byte(md.UID), b); err != nil {
			panic(err)
		}

		/*
			for i := 0; i < 1000; i++ {

				uid = uuid.New()
				email := fake.Email()
				md.Email = email
				md.Pass = fmt.Sprintf("%x", md5.Sum([]byte("1234567890")))
				md.UID = uid.String()
				md.Root = false

				b, err = md.MarshalBinary()
				if err != nil {
					panic(err)
				}

				if err = bk.Put([]byte(email), b); err != nil {
					panic(err)
				}
			}
		*/
		return nil
	})

	assert.NoError(t, err)
}
