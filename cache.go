package main

import (
	"fmt"
	"os"
	"github.com/boltdb/bolt"
)

type CacheDB struct {
	database *bolt.DB
	bucket string
}

func NewCacheDb(filename string, bucketname string) *CacheDB {
	db := CacheDB{}
	var database *bolt.DB
	
	// Create/open the database
	
	database, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		fmt.Sprintf("Error, cannot open the database %v", filename)
		fmt.Println(err)
		os.Exit(1)
	}

	// Create the bucket if it doesn't exist
	database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("RFIDBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	
	if err != nil {
		fmt.Sprintf("Error, cannot create or open the bucket %v", bucketname)
		fmt.Println(err)
		os.Exit(1)
	}

	
	db.database = database
	db.bucket = bucketname
	return &db
}

func (db *CacheDB) Close() error {
	return db.database.Close()
}

func (db *CacheDB) checkCacheDBForTag(tag string) bool {
	val := ""
	db.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		val = string(b.Get([]byte(tag)))
		return nil
	})

	if val != "" {
		return true
	}
	return false
}

func (db *CacheDB) addTagToCacheDB(tag string) {
	db.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte(tag), []byte(tag))
		return err
	})
}
