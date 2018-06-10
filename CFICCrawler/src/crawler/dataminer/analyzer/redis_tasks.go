package analyzer

import (
	"routingpool"
	"htmlparser"
	"github.com/spf13/viper"
	"time"
	"os"
	"bufio"
	"github.com/axgle/mahonia"
	"strings"
	"fmt"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPush *RedisPushTask
	RedisChange *RedisChangeTask

	RedisConnection redis.Conn
)

func init() {
	RedisConnection = initRedis()

	RedisPush = NewRedisPushTask()
	RedisChange = NewRedisChangeTask()


}

func initRedis() redis.Conn{
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logger.Errorf("Connect to redis error", err)
	}

	return c
}

type RedisChangeTask struct {
	*routingpool.Base

	pushDone chan bool
	redis redis.Conn
}

// Implement Task interface
type RedisPushTask struct {
	*routingpool.Base
	message chan interface{}

	redis redis.Conn
	//timer int	// The waiting seconds for receiving data, analysis routine exits after the waiting.
}

func NewRedisChangeTask() *RedisChangeTask {
	return &RedisChangeTask{pushDone:make(chan bool), Base: &routingpool.Base{Name: "Redis Change Task", Response: make(chan bool)}, redis:RedisConnection}
}

func (r *RedisChangeTask) Run(id int) {
	r.caller(id)
}

func (r *RedisChangeTask) caller(id int) {
	count := 0
	exit := false

	for !exit {
		select {
		case <-r.pushDone:
			count = count + 1
			if count == viper.GetInt("redis.pushtask.count") {
				exit = true
			}
		}
	}

	var data map[string]string
	encoder := mahonia.NewEncoder("gbk")
	funds, _ := getFunds()
	for _,f := range funds{
		logger.Debugf("Trying to get records from Redis for %s.", f)
		key := encoder.ConvertString(f)
		values, err := redis.Values(r.redis.Do("LRANGE", key, -2, -1))
		if err != nil{
			logger.Errorf("redis lrange failed:", err)
		}

		for _, v := range values{
			json.Unmarshal(v.([]byte), &data)
			logger.Infof("Found record - code:%s, recorddata:%s, holdcount::%s, holdvalue:%s", data["code"], data["recorddate"], data["holdcount"], data["holdvalue"])
		}
	}

	if r.redis != nil {
		r.redis.Close()
	}
}

func NewRedisPushTask() *RedisPushTask {
	return &RedisPushTask{message : make(chan interface{}, 1024), Base:&routingpool.Base{Name: "Redis Push Task", Response: make(chan bool)}, redis:RedisConnection}
}

func PushDataIntoRedis(msg interface{}) {
	RedisPush.message <- msg
}

func PushTaskDone() {
	RedisChange.pushDone <- true
}

func (a *RedisPushTask) caller(id int) {
	timeout := time.NewTimer(time.Second * time.Duration(viper.GetInt("redis.pushtask.timer")))
	exit := false
	funds, _ := getFunds()
	encoder := mahonia.NewEncoder("gbk")

	for !exit {
		select {
			case data := <-a.message:
				logger.Infof("Analysis-task %d, received data", id)
				tmp := data.([]*htmlparser.JJCCData)

				for _, value := range tmp {
					logger.Infof("Analysis-task %d, row data name %s, code %s, holdcount %.4f, holdvalue %.4f", id, value.Name, value.Code, value.HoldCount, value.HoldValue)
					for _,fund := range funds {
						if strings.Contains(value.Name, fund) {
							logger.Debug("Analysis-task %d, found record data for %s", value.Name)

							raw := map[string]string{	"code" : value.Code,
														"recorddate" : value.RecordDate,
														"holdcount" : fmt.Sprintf("%.4f", value.HoldCount),
														"holdvalue" : fmt.Sprintf("%.4f", value.HoldValue),
													}

							key := encoder.ConvertString(fund)
							json_value, _ := json.Marshal(raw)
							_, err := a.redis.Do("LPUSH", key, json_value)
							if err != nil {
								logger.Errorf("redis set failed:", err)
							}
						}
					}
				}

				timeout.Reset(time.Second * time.Duration(viper.GetInt("analyser.timer")))

			case <- timeout.C: // The waiting seconds for receiving data, analysis routine exits after the waiting.
				PushTaskDone()
				exit = true
				break
		}
	}

	logger.Infof("Analysis-task %d, exit.....................", id)
}

func (a *RedisPushTask) Run(id int) {
	a.caller(id)
}

func getFunds() ([]string, error) {
	filename := viper.GetString("module.jjcc.funds_file_path")
	file, err:= os.Open(filename)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "\n", filename, err)
		logger.Errorf("WARN: Open file %s failed, %s", filename, err)
		return nil, err
	}
	defer file.Close()

	var result []string
	decoder := mahonia.NewDecoder("gbk")
	scanner := bufio.NewScanner(decoder.NewReader(file))
	for scanner.Scan() {
		//fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
		result = append(result, scanner.Text())
	}

	return result,nil
}