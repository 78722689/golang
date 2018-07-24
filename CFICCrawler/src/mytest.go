package main

import "fmt"

func redis_practice() {
      /* c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("Connect to redis error", err)
			return
		}
		defer c.Close()
      */

      fmt.Println(401/10)
      length :=400/10
      if 400%10 !=0 {
      	length = 401/10+1
	  }
	  x:=[]string{"a","b", "c","d","e","f"}
      for index :=0; index<=length;index++ {
      	fmt.Println(x[index:index])
	  }

	/*4
		imap := mastring]string{"username": "666", "phonenumber": "888"}
		value, _ := json.Marshal(imap)

		_, err = c.Do("LPUSH", "mylist", value)
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
	*/
/*
    mylist, err := redis.Strings(c.Do("LRANGE", "招商中证银行指数分级证券投资基金", 0, 9))
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

func main() {
    redis_practice()

}
