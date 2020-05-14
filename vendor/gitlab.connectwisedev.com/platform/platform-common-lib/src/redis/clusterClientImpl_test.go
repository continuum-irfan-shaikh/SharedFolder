package redis

//TODO : https://continuum.atlassian.net/browse/RMM-58213

// import (
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/alicebob/miniredis"
// 	"github.com/go-redis/redis"
// )

// func newTestClusterRedis() *redis.ClusterClient {
// 	mr, err := miniredis.Run()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// client := redis.NewClient(&redis.Options{
// 	// 	Addr: mr.Addr(),
// 	// })
// 	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs: []string{mr.Addr()},
// 	})

// 	return clusterClient
// }
// func Test_redisClusterFunctionality(t *testing.T) {
// 	z := Z{Score: 0, Member: "abc"}
// 	z1 := Z{Score: 0, Member: "xyz"}

// 	t.Run("Add_member_to_a_sorted_set,_or_update_its_score_if_it_already_exists", func(t *testing.T) {

// 		c := &clusterClientImpl{
// 			config:        &Config{},
// 			clusterClient: newTestClusterRedis(),
// 		}

// 		c.ZAdd("k1", z, z1)
// 		strarray, err := c.ZRange("k1", 0, -1)
// 		if err != nil {
// 			t.Errorf("clientImpl.ZAdd() error = %v, wantErr %s", err, "nil")
// 			return
// 		}
// 		if strarray[0] != z.Member || strarray[1] != z1.Member {
// 			t.Errorf("clientImpl.ZAdd() = %v and %v, want %v and %v", z.Member, z1.Member, strarray[0], strarray[1])
// 		}
// 		_, err = c.ZRem("k1", z1.Member)

// 		if err != nil {
// 			t.Errorf("error of ZRem() while removing member; error = %v, wantErr %s", err, "nil")
// 			return
// 		}
// 		strarray, err = c.ZRange("k1", 0, -1)
// 		if strarray[0] != z.Member || len(strarray) != 1 {
// 			t.Errorf("after removal of element length of array is not equal to 1 or got value is %v but want value %v", strarray[0], z.Member)
// 		}
// 		output, err := c.Exists("k1")
// 		if err != nil {
// 			t.Errorf("error of ZRem() while removing member; error = %v, wantErr %s", err, "nil")
// 			return
// 		}
// 		if output != 1 {
// 			t.Errorf("No of element is set with k1 key is equal to: %v but no of elements want: 1", output)
// 		}
// 		c.Close()
// 	})
// }

// func TestGetClusterClientService(t *testing.T) {
// 	conf := &Config{}
// 	type args struct {
// 		transactionID string
// 		config        *Config
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want Client
// 	}{
// 		{
// 			name: "nil_config",
// 			args: args{
// 				transactionID: "abc",
// 				config:        conf,
// 			},
// 			want: &clusterClientImpl{config: conf},
// 		},
// 		{
// 			name: "empty_transaction_id",
// 			args: args{
// 				transactionID: "",
// 				config:        conf,
// 			},
// 			want: &clusterClientImpl{config: conf},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetClusterClientService(tt.args.transactionID, tt.args.config); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetClusterClientService() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_clusterClientImpl_Init(t *testing.T) {
// 	type fields struct {
// 		config        *Config
// 		clusterClient *redis.ClusterClient
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			name: "cluster_client_not_nil",
// 			fields: fields{
// 				config:        &Config{},
// 				clusterClient: &redis.ClusterClient{},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "cluster_client_nil",
// 			fields: fields{
// 				config:        &Config{},
// 				clusterClient: nil,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "cluster_client_config_nil",
// 			fields: fields{
// 				config:        nil,
// 				clusterClient: nil,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &clusterClientImpl{
// 				config:        tt.fields.config,
// 				clusterClient: tt.fields.clusterClient,
// 			}
// 			if err := c.Init(); (err != nil) != tt.wantErr {
// 				t.Errorf("clusterClientImpl.Init() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_clusterClientImpl_Close(t *testing.T) {
// 	type fields struct {
// 		config        *Config
// 		clusterClient *redis.ClusterClient
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			name: "cluster_client_not_nil",
// 			fields: fields{
// 				config:        &Config{},
// 				clusterClient: nil,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &clusterClientImpl{
// 				config:        tt.fields.config,
// 				clusterClient: tt.fields.clusterClient,
// 			}
// 			if err := c.Close(); (err != nil) != tt.wantErr {
// 				t.Errorf("clusterClientImpl.Close() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_stringRedisClusterFunctionality(t *testing.T) {

// 	t.Run("Add_member_to_a_sorted_set,_or_update_its_score_if_it_already_exists", func(t *testing.T) {

// 		c := &clusterClientImpl{
// 			config:        &Config{},
// 			clusterClient: newTestClusterRedis(),
// 		}

// 		c.Set("k1", "abc")
// 		strarray, err := c.Get("k1")
// 		if err != nil {
// 			t.Errorf("clientImpl.ZAdd() error = %v, wantErr %s", err, "nil")
// 			return
// 		}

// 		err = c.Delete("k1")

// 		if err != nil {
// 			t.Errorf("error of ZRem() while removing member; error = %v, wantErr %s", err, "nil")
// 			return
// 		}
// 		strarray, err = c.Get("k1")
// 		if strarray != "" {
// 			t.Errorf("after removal of element length of array is not equal to 1 or got value is %v but want value %v", strarray, "")
// 		}
// 		c.SetWithExpiry("k1", "xyz", 1000)
// 		time.Sleep(time.Second)
// 		strarray, err = c.Get("k1")

// 		if strarray != "" {
// 			t.Errorf("after removal of element length of array is not equal to 1 or got value is %v but want value %v", strarray, "abc")
// 		}

// 		c.Set("k1", "xyz")

// 		c.Expire("k1", 1)
// 		time.Sleep(time.Second)
// 		strarray, err = c.Get("k1")
// 		if strarray != "" {
// 			t.Errorf("after removal of element length of array is not equal to 1 or got value is %v but want value %v", strarray, "abc")
// 		}
// 		var value int64
// 		value = 2
// 		c.Set("k1", value)
// 		p, _ := c.Incr("k1")
// 		if (value + 1) != (p) {
// 			t.Errorf("Expected %v and Got %v", p, (value + 1))

// 		}
// 		c.Delete("k1")
// 		c.Set("k1", "abc")
// 		c.Set("k2", "xyz")
// 		c.Set("k3", "pqr")
// 		var cursor uint64
// 		keys, cursor, _ := c.Scan(cursor, "k*", 1)
// 		if len(keys) != 3 {
// 			t.Errorf("Expected 3 and Got %v", len(keys))
// 		}

// 		c.CreatePipeline()

// 		c.Close()
// 	})
// }

// func Test_clusterClientImpl_Incr(t *testing.T) {
// 	type fields struct {
// 		config        *Config
// 		clusterClient *redis.ClusterClient
// 	}
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    int64
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &clusterClientImpl{
// 				config:        tt.fields.config,
// 				clusterClient: tt.fields.clusterClient,
// 			}
// 			got, err := c.Incr(tt.args.key)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("clusterClientImpl.Incr() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("clusterClientImpl.Incr() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_clusterClientImpl_SetWithExpiry(t *testing.T) {
// 	type fields struct {
// 		config        *Config
// 		clusterClient *redis.ClusterClient
// 	}
// 	type args struct {
// 		key      string
// 		value    interface{}
// 		duration time.Duration
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &clusterClientImpl{
// 				config:        tt.fields.config,
// 				clusterClient: tt.fields.clusterClient,
// 			}
// 			if err := c.SetWithExpiry(tt.args.key, tt.args.value, tt.args.duration); (err != nil) != tt.wantErr {
// 				t.Errorf("clusterClientImpl.SetWithExpiry() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
