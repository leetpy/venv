package manager

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	. "github.com/leetpy/venv/internal"
)

type dbAdapter interface {
	Init() error
	//ListServer()
	//UpdateServer()
	//AddServer()
	//DelServer()
	ListTypes() ([]NodeType, error)
	CreateType(t NodeType) (NodeType, error)
	//Close() error
}

func makeDBAdapter(dbType string, dbFile string) (dbAdapter, error) {
	if dbType == "bolt" {
		innerDB, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			return nil, err
		}
		db := boltAdapter{
			db:     innerDB,
			dbFile: dbFile,
		}
		err = db.Init()
		return &db, err
	}
	// unsupported db-type
	return nil, fmt.Errorf("unsupported db-type: %s", dbType)
}

const (
	_nodeBucketKey = "node"
	_portBucketKey = "port"
	_typeBucketKey = "type"
)

type boltAdapter struct {
	db     *bolt.DB
	dbFile string
}

func (b *boltAdapter) Init() (err error) {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(_nodeBucketKey))
		if err != nil {
			return fmt.Errorf("create bucket %s error: %s", _nodeBucketKey, err.Error())
		}

		_, err = tx.CreateBucketIfNotExists([]byte(_portBucketKey))
		if err != nil {
			return fmt.Errorf("create bucket %s error: %s", _portBucketKey, err.Error())
		}
		_, err = tx.CreateBucketIfNotExists([]byte(_typeBucketKey))
		if err != nil {
			return fmt.Errorf("create bucket %s error: %s", _typeBucketKey, err.Error())
		}
		return nil
	})
}

func (b *boltAdapter) ListTypes() (nt []NodeType, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_typeBucketKey))
		c := bucket.Cursor()
		var t NodeType
		for k, v := c.First(); k != nil; k, v = c.Next() {
			jsonErr := json.Unmarshal(v, &t)
			if jsonErr != nil {
				err = fmt.Errorf("%s; %s", err.Error(), jsonErr)
				continue
			}
			nt = append(nt, t)
		}
		return err
	})
	return
}

func (b *boltAdapter) CreateType(t NodeType) (NodeType, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_typeBucketKey))
		v, err := json.Marshal(t)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(t.ID), v)
		return err
	})
	return t, err
}
