package heartbeat

import (
	"fmt"
	"time"
	"strconv"

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

func SetHeartbeat(uid string, heartbeat model.HeartbeatStatus) error {
	heartbeatTime := time.Now()

	// By all accounts, for the record, err should always be nil
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(uid + "-status"), []byte(heartbeat))
		if err != nil {
			return err
		}

		err = txn.Set([]byte(uid + "-time"), []byte(fmt.Sprintf("%d", heartbeatTime.Unix())))
		return err
	})

	return err
}

func GetHeartbeat(uid string) (model.Heartbeat, error) {
	var heartbeat model.Heartbeat
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(uid + "-status"))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}

		status, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte(uid + "-time"))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}

		lastSeenRaw, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		lastSeen, err := strconv.ParseInt(string(lastSeenRaw), 10, 64)
		if err != nil {
			return err
		}

		heartbeat = model.Heartbeat {
			Status: model.HeartbeatStatus(status),
			LastSeen: int(lastSeen),
		}

		return err
	})

	return heartbeat, err
}
