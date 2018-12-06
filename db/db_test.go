package db

import (
	"context"
	"os"
	"testing"

	"github.com/pingcap/tidb/store/mockstore"
	"github.com/stretchr/testify/assert"
)

var mockDB *DB

func TestMain(m *testing.M) {
	store, err := mockstore.NewMockTikvStore()
	if err != nil {
		panic(err)
	}
	mockDB = &DB{
		Namespace: "mockdb-ns",
		ID:        1,
		kv:        &RedisStore{store},
	}

	os.Exit(m.Run())
}

func MockTest(t *testing.T, callFunc func(txn *Transaction)) {
	txn, err := mockDB.Begin()
	assert.NoError(t, err)
	callFunc(txn)
	err = txn.Commit(context.TODO())
	assert.NoError(t, err)
}
