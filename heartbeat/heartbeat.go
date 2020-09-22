package heartbeat

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/solderneer/axiom-backend/graph/model"
)

var db *badger.DB

func Close() {
	db.Close()
}

func InitHeartbeat(badgerDir string) error {
	var err error
	db, err = openBadger(badgerDir)

	return err
}

func openBadger(badgerDir string) (*badger.DB, error) {
	if badgerDir == "" {
		return badger.Open(badger.DefaultOptions("").WithInMemory(true))
	} else {
		return badger.Open(badger.DefaultOptions(badgerDir))
	}
}

func SetHeartbeat(uid string, heartbeat model.Heartbeat) error {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(uid), []byte(heartbeat))
	})

	return err
}

func GetHeartbeat(uid string) (model.Heartbeat, error) {
	var heartbeat model.Heartbeat
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uid))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		heartbeat = model.Heartbeat(val)

		return err
	})

	return heartbeat, err
}
