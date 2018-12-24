package mpc

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	MPC_TASK_KEY_ALICE  = "MPC_TASK_LIST_ALICE_KEY"
	MPC_TASK_KEY_BOB = "MPC_TASK_LIST_BOB_KEY"
)

type Redis struct {
	Con redis.Conn
}

func NewRedis(url string) (*Redis, error) {
	con, err := redis.Dial("tcp", url)
	if err != nil {
		return &Redis{}, err
	}
	return &Redis{
		Con: con,
	}, nil
}

func (r *Redis) RPush(key string, elem map[string]string) error {
	jsonStr, _ := json.Marshal(elem)
	_, err := r.Con.Do("rpush", key, string(jsonStr))
	if err != nil {
		fmt.Println("rpush set failed:", err.Error())
		return err
	}
	return nil
}

func (r *Redis) Values(key string) []map[string]string {
	values, _ := redis.Values(r.Con.Do("lrange", key, "0", "100"))
	res := make([]map[string]string, 0)
	for _, v := range values {
		var tmpMap map[string]string
		if byt, ok := v.([]byte); ok {
			json.Unmarshal(byt, &tmpMap)
			res = append(res, tmpMap)
		}
	}
	return res
}

/// return data of top.
func (r *Redis) RPop(key string) string {
	res, err := redis.String(r.Con.Do("lpop", key))
	if err != nil {
		return ""
	}
	return res
}

