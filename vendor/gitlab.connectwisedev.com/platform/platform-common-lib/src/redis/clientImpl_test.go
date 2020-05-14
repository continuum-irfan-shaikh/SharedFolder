package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func newTestRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

func Test_redisFunctionality(t *testing.T) {
	z := Z{Score: 0, Member: "abc"}
	z1 := Z{Score: 0, Member: "xyz"}

	t.Run("Add_member_to_a_sorted_set,_or_update_its_score_if_it_already_exists", func(t *testing.T) {

		c := &clientImpl{
			config: &Config{},
			client: newTestRedis(),
		}
		c.ZAdd("k1", z, z1)
		strarray, err := c.ZRange("k1", 0, -1)
		if err != nil {
			t.Errorf("clientImpl.ZAdd() error = %v, wantErr %s", err, "nil")
			return
		}
		if strarray[0] != z.Member || strarray[1] != z1.Member {
			t.Errorf("clientImpl.ZAdd() = %v and %v, want %v and %v", z.Member, z1.Member, strarray[0], strarray[1])
		}
		_, err = c.ZRem("k1", z1.Member)

		if err != nil {
			t.Errorf("error of ZRem() while removing member; error = %v, wantErr %s", err, "nil")
			return
		}
		strarray, err = c.ZRange("k1", 0, -1)
		if strarray[0] != z.Member || len(strarray) != 1 {
			t.Errorf("after removal of element length of array is not equal to 1 or got value is %v but want value %v", strarray[0], z.Member)
		}
		output, err := c.Exists("k1")
		if err != nil {
			t.Errorf("error of ZRem() while removing member; error = %v, wantErr %s", err, "nil")
			return
		}
		if output != 1 {
			t.Errorf("No of element is set with k1 key is equal to: %v but no of elements want: 1", output)
		}
	})
}

func TestIncr(t *testing.T) {
	redis := newTestRedis()
	client := clientImpl{client: redis}

	client.Incr("counter")

	counter, _ := redis.Get("counter").Result()
	if counter != "1" {
		t.Errorf("expected: 1, got: %v", counter)
	}
}

func TestExpire(t *testing.T) {
	redis := newTestRedis()
	client := clientImpl{client: redis}

	redis.Set("cache", "data", 0)

	client.Expire("cache", 10*time.Second)

	ttl, _ := redis.TTL("cache").Result()
	if ttl.Seconds() == -1 {
		t.Errorf("key doesn't have associated expire")
	}
	if ttl.Seconds() == -2 {
		t.Errorf("key not found")
	}
}
