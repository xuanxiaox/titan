package command

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initHashes(t *testing.T, key string, n int) {
	args := []string{key}
	for i := n; i > 0; i-- {
		args = append(args, strconv.Itoa(i), "bar")
	}
	ctx := ContextTest("hmset", args...)
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, "+OK", lines[0])
}

func clearHashes(t *testing.T, key string) {
	ctx := ContextTest("del", key)
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])
}

func setHashes(t *testing.T, args ...string) []string {
	ctx := ContextTest("hmset", args...)
	Call(ctx)
	return ctxLines(ctx.Out)
}

func TestHLen(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hlen", key)
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, ":3", lines[0])

	// case 2
	lines = setHashes(t, key, "a", "a", "b", "b")
	assert.Equal(t, "+OK", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":5", lines[0])

	// case 3
	lines = setHashes(t, key, "c", "c", "c", "d")
	assert.Equal(t, "+OK", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":6", lines[0])

	// end
	clearHashes(t, key)
}

func TestHDel(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 5)

	// case 1
	ctx := ContextTest("hdel", key, "1")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":4", lines[0])

	// case 2
	ctx = ContextTest("hdel", key, "2")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":3", lines[0])

	// case 3
	ctx = ContextTest("hdel", key, "3", "4", "5")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":3", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":0", lines[0])
	// then re-insert into hash
	lines = setHashes(t, key, "a", "a", "b", "b")
	assert.Equal(t, "+OK", lines[0])
	ctx = ContextTest("hlen", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":2", lines[0])
}

func TestHSetNX(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hsetnx", key, "1", "haha")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, ":0", lines[0])
	ctx = ContextTest("hget", key, "1")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$3", lines[0])
	assert.Equal(t, "bar", lines[1])

	// case 2
	ctx = ContextTest("hsetnx", key, "4", "haha")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])
	ctx = ContextTest("hget", key, "4")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$4", lines[0])
	assert.Equal(t, "haha", lines[1])

	// case 3
	ctx = ContextTest("hsetnx", key, "4", "bar")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":0", lines[0])
	ctx = ContextTest("hget", key, "4")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$4", lines[0])
	assert.Equal(t, "haha", lines[1])

	// end
	clearHashes(t, key)
}
func TestHSet(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hset", key, "1", "haha")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, ":0", lines[0])
	ctx = ContextTest("hget", key, "1")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$4", lines[0])
	assert.Equal(t, "haha", lines[1])

	// case 2
	ctx = ContextTest("hset", key, "4", "haha")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])
	ctx = ContextTest("hget", key, "4")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$4", lines[0])
	assert.Equal(t, "haha", lines[1])

	// case 3
	ctx = ContextTest("hset", key, "4", "bar")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":0", lines[0])
	ctx = ContextTest("hget", key, "4")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$3", lines[0])
	assert.Equal(t, "bar", lines[1])

	// end
	clearHashes(t, key)
}

func TestHGet(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hget", key, "1")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, "$3", lines[0])
	assert.Equal(t, "bar", lines[1])

	// case 2
	ctx = ContextTest("hget", key, "5")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "$-1", lines[0])

	// end
	clearHashes(t, key)
}

func TestHGetAll(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hgetall", key)
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, "*6", lines[0])
	assert.Equal(t, "$1", lines[1])
	assert.Equal(t, "1", lines[2])
	assert.Equal(t, "$3", lines[3])
	assert.Equal(t, "bar", lines[4])
	assert.Equal(t, "$1", lines[5])
	assert.Equal(t, "2", lines[6])
	assert.Equal(t, "$3", lines[7])
	assert.Equal(t, "bar", lines[8])
	assert.Equal(t, "$1", lines[9])
	assert.Equal(t, "3", lines[10])
	assert.Equal(t, "$3", lines[11])
	assert.Equal(t, "bar", lines[12])

	// case 2
	ctx = ContextTest("hset", key, "foo", "haha")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])

	ctx = ContextTest("hgetall", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "*8", lines[0])
	assert.Equal(t, "$1", lines[1])
	assert.Equal(t, "1", lines[2])
	assert.Equal(t, "$3", lines[3])
	assert.Equal(t, "bar", lines[4])
	assert.Equal(t, "$1", lines[5])
	assert.Equal(t, "2", lines[6])
	assert.Equal(t, "$3", lines[7])
	assert.Equal(t, "bar", lines[8])
	assert.Equal(t, "$1", lines[9])
	assert.Equal(t, "3", lines[10])
	assert.Equal(t, "$3", lines[11])
	assert.Equal(t, "bar", lines[12])
	assert.Equal(t, "$3", lines[13])
	assert.Equal(t, "foo", lines[14])
	assert.Equal(t, "$4", lines[15])
	assert.Equal(t, "haha", lines[16])

	// end
	clearHashes(t, key)
}

func TestHMGet(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hmget", key, "1", "2", "3")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, "*3", lines[0])
	assert.Equal(t, "$3", lines[1])
	assert.Equal(t, "bar", lines[2])
	assert.Equal(t, "$3", lines[3])
	assert.Equal(t, "bar", lines[4])
	assert.Equal(t, "$3", lines[5])
	assert.Equal(t, "bar", lines[6])

	// case 2
	ctx = ContextTest("hset", key, "foo", "haha")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, ":1", lines[0])

	ctx = ContextTest("hmget", key, "1", "foo")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "*2", lines[0])
	assert.Equal(t, "$3", lines[1])
	assert.Equal(t, "bar", lines[2])
	assert.Equal(t, "$4", lines[3])
	assert.Equal(t, "haha", lines[4])

	// end
	clearHashes(t, key)
}
func TestHMSet(t *testing.T) {
	// init
	key := "hash-key"
	initHashes(t, key, 3)

	// case 1
	ctx := ContextTest("hmset", key, "1", "ha", "2", "haha", "3", "hahaha")
	Call(ctx)
	lines := ctxLines(ctx.Out)
	assert.Equal(t, "+OK", lines[0])
	ctx = ContextTest("hgetall", key)
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "*6", lines[0])
	assert.Equal(t, "1", lines[2])
	assert.Equal(t, "ha", lines[4])
	assert.Equal(t, "2", lines[6])
	assert.Equal(t, "haha", lines[8])
	assert.Equal(t, "3", lines[10])
	assert.Equal(t, "hahaha", lines[12])

	// case 2
	ctx = ContextTest("hmset", key, "foo", "bar")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "+OK", lines[0])
	ctx = ContextTest("hmget", key, "1", "foo")
	Call(ctx)
	lines = ctxLines(ctx.Out)
	assert.Equal(t, "*2", lines[0])
	assert.Equal(t, "$2", lines[1])
	assert.Equal(t, "ha", lines[2])
	assert.Equal(t, "$3", lines[3])
	assert.Equal(t, "bar", lines[4])

	// end
	clearHashes(t, key)
}

func TestHExists(t *testing.T)      {}
func TestHIncrBy(t *testing.T)      {}
func TestHIncrByFloat(t *testing.T) {}
func TestHKeys(t *testing.T)        {}
func TestHStrLen(t *testing.T)      {}
func TestHVals(t *testing.T)        {}
func TestHMSlot(t *testing.T)       {}