package mpc

import (
	"testing"
)

func Test_Put2List(t *testing.T) {
	redis, err := NewRedis("192.168.9.14:6379")
	if err != nil {
		t.Fatal("new redis instance fail: ", err.Error())
	}
	v := make(map[string]string)
	v["taskId"] = "0xxxxxxxxxxxxxxxxxxxxxxxxxxx"
	v["pubKey"] = "0xxxlllll"
	err = redis.RPush(MPC_TASK_KEY_ALICE, v)
	err = redis.RPush(MPC_TASK_KEY_BOB, v)

	v["taskId"] = "0xxxxxxxxxxxxxxxxxxxxxxxxxxx"
	v["pubKey"] = "0xxx222"
	err = redis.RPush(MPC_TASK_KEY_ALICE, v)
	err = redis.RPush(MPC_TASK_KEY_BOB, v)

	if err != nil {
		t.Fatal("lpush set fail : ", err.Error())
	}
	// query result
	//strs := redis.Values()
	/*for _, v := range strs {
		for key, val := range v {
			t.Logf("key: %v, value: %v", key, val)
		}
	}*/

	/*for i := 0; i < len(strs); i++ {
		t.Log(redis.RPop())
	}*/
}