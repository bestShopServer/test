package boltdb

import (
	"DetectiveMasterServer/util"
	"fmt"
	"github.com/boltdb/bolt"
)

const DB_CATEGORY = "Database"

var DB *bolt.DB

// Func: Create Database If Not Exist
func CreateDatabase() *bolt.DB {
	fmt.Println("CreateDatabase ...")
	db, err := bolt.Open("game.db", 0600, nil)
	fmt.Println("err:", err)
	if err != nil {
		util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Open Database Err:"+err.Error())
	}
	fmt.Println("CreateDatabase db:", db)
	return db
}

// Func: Create Bucket in Database
func CreateBucket(name string) {
	err := DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Create "+name+" Bucket Err:"+err.Error())
		}
		return err
	})
	if err != nil {
		util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Update Database Err:"+err.Error())
	}
}

// Func: Create Or Update Key Value In Bucket
func CreateOrUpdate(key []byte, value []byte, bucket string) {
	err := DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
	if err != nil {
		util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Update Database Err:"+err.Error())
	}
}

// Func: View Key Value In Bucket
func View(key []byte, bucket string) []byte {
	var ret []byte
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		ret = b.Get(key)
		return nil
	})
	return ret
}

// Func: Delete Bucket By Name
func DeleteBucket(bucket string) {
	err := DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
	if err != nil {
		util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Update Database Err:"+err.Error())
	}
}

// Func: Delete Key Value From Bucket
func Delete(key []byte, bucket string) {
	err := DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(key)
		return err
	})
	if err != nil {
		util.Logger(util.ERROR_LEVEL, DB_CATEGORY, "Update Database Err:"+err.Error())
	}
}
