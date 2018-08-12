package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/axgle/mahonia"
	"time"
	"flag"
	"encoding/json"
)

func redis_practice() {
       c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("Connect to redis error", err)
			return
		}
		defer c.Close()
	encoder := mahonia.NewEncoder("gbk")
	key := encoder.ConvertString("上证大宗商品股票交易型开放式指数证券投资基金")
		_, err = c.Do("DEL", key)
		fmt.Println(err)

	/*4
		imap := mastring]string{"username": "666", "phonenumber": "888"}
		value, _ := json.Marshal(imap)

		_, err = c.Do("LPUSH", "mylist", value)
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
	*/

/*
    mylist, err := redis.Strings(c.Do("LRANGE", key, 0, 9))
    if err != nil {
        fmt.Println("redis get failed:", err)
    } else {
        for _,value := range mylist{
            var imapGet map[string]string
            errShal := json.Unmarshal([]byte(value), &imapGet)
            if errShal != nil {
                fmt.Println(err)
            }

            fmt.Println(imapGet["code"], imapGet["recorddate"], imapGet["holdcount"], imapGet["holdvalue"])
        }
    }
*/
/*
    encoder := mahonia.NewEncoder("gbk")
    key := encoder.ConvertString("上证大宗商品股票交易型开放式指数证券投资基金")
    //value := encoder.ConvertString("蚊子-z")
    result, err := redis.Values(c.Do("LRANGE", key , -2, -1))
    if err != nil {
        fmt.Println("redis get failed:", err)
    }
    fmt.Println(result)
    var m map[string]string
	for _, v:= range result {
		json.Unmarshal(v.([]byte), &m)
		fmt.Println(m)
		//fmt.Println(m["recorddate"])
	}

*/
}


func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

var (
	pool *redis.Pool
	redisServer = flag.String("127.0.0.1", ":6379", "")
)

func main() {
    //redis_practice()
	flag.Parse()
	pool = newPool(*redisServer)

	s,_ := pool.Get().Do("get", "DOMAIN_600007")
	//for _,value := range mylist {
		var imapGet []string
		json.Unmarshal(s.([]byte), &imapGet)

		fmt.Println(fmt.Sprintf("%v", imapGet))
	//}
	//mymap := make(map[string]map[string]map[string]string)
/*
	temp2 := map[string]map[string]string{}
	for index := 0;index != 5; index++ {
		temp1 := map[string]string{"count": fmt.Sprintf("8%d", index), "value": fmt.Sprintf("9%d", index)}
		temp2["20120" + fmt.Sprintf("%d", index)] = temp1
	}
	temp3 := map[string]map[string]map[string]string{"2018-08-31":temp2}

	json_value, _ := json.Marshal(temp3)
	pool.Get().Do("LPUSH", "test", json_value)
	fmt.Println(temp3["2018-08-31"]["201202"])

	encoder := mahonia.NewEncoder("gbk")
	fund_info := encoder.ConvertString("1" + "_" +"2")
	_, err := pool.Get().Do("SADD", "FUND_INFO_TABLE", fund_info)
	if err != nil {
		fmt.Println("Push fund info to Redis failure:", err)
	}
*/
}
