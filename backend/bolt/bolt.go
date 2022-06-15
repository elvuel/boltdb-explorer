package bolt

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

func Read(path string) (map[string]interface{}, error) {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result, err := ToMap(db)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func WriteMap(path string, data map[string]interface{}) (string, error) {
	db, err := FromMap(data, path)
	if err != nil {
		os.Remove(path)
		return path, err
	}
	db.Close()

	return path, nil
}

func Write(path string, data []byte) (string, error) {
	o := make(map[string]interface{})
	err := json.Unmarshal(data, &o)
	if err != nil {
		return "", err
	}

	db, err := FromMap(o, path)
	if err != nil {
		os.Remove(path)
		return path, err
	}
	db.Close()

	return path, nil
}

func FromMap(data map[string]interface{}, store string) (*bolt.DB, error) {
	db, err := bolt.Open(store, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for k, v := range data {
			b, _ := tx.CreateBucketIfNotExists([]byte(k))
			err := RecursiveToDB(b, k, v)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return db, err
}

func RecursiveToDB(b *bolt.Bucket, k string, v interface{}) error {
	switch v.(type) {
	case map[string]interface{}:
		for k1, v1 := range v.(map[string]interface{}) {
			if _, ok := v1.(map[string]interface{}); ok {
				b1, _ := b.CreateBucketIfNotExists([]byte(k1))
				err := RecursiveToDB(b1, k1, v1)
				if err != nil {
					return err
				}
			} else {
				err := RecursiveToDB(b, k1, v1)
				if err != nil {
					return err
				}
			}
		}
	case string, int, bool, float64, float32, int64, int32, int16, int8, uint64, uint32, uint16, uint8:
		err := b.Put([]byte(k), []byte(fmt.Sprintf("%v", v)))
		if err != nil {
			return err
		}
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = b.Put([]byte(k), data)
		if err != nil {
			return err
		}
	}
	return nil
}

func ToMap(db *bolt.DB) (result map[string]interface{}, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		o := make(map[string]interface{})
		result = RecursiveToMap(tx, c, o)
		return nil
	})
	return
}

func RecursiveToMap(tx *bolt.Tx, c *bolt.Cursor, o map[string]interface{}) map[string]interface{} {
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			bucket := c.Bucket().Bucket(k)
			if bucket == nil {
				bucket = tx.Bucket(k)
			}
			nc := bucket.Cursor()
			no := make(map[string]interface{})
			o[string(k)] = RecursiveToMap(tx, nc, no)
			continue
		}
		o[string(k)] = string(v)
	}
	return o
}
