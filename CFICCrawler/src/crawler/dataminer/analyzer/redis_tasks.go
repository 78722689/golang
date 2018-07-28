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
	"utility"
)

var (
	RedisPush *RedisPushTask
	RedisDispatcher *RedisManager

	RedisConnection redis.Conn
	redisLangEncoder mahonia.Encoder
)

func init() {
	RedisConnection = initRedis()
	redisLangEncoder = mahonia.NewEncoder("gbk")

	RedisPush = NewRedisPushTask()
	RedisDispatcher = NewRedisManagerTask()
}

func initRedis() redis.Conn{
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logger.Errorf("Connect to redis error", err)
	}

	return c
}

type RedisManager struct {
	*routingpool.Base

	pushDone chan bool
	changeDone chan bool
	redis redis.Conn
}

// Implement Task interface
type RedisPushTask struct {
	*routingpool.Base
	message chan interface{}

	encoder mahonia.Encoder
	redis redis.Conn
	//timer int	// The waiting seconds for receiving data, analysis routine exits after the waiting.
}

// Implement Task interface
type RedisChangeTask struct {
	*routingpool.Base
	funds []string

	encoder mahonia.Encoder
	redis redis.Conn
}

func NewRedisManagerTask() *RedisManager {
	return &RedisManager{pushDone:make(chan bool), changeDone:make(chan bool), Base: &routingpool.Base{Name: "Redis Change Task", Response: make(chan bool)}, redis:RedisConnection}
}

func NewRedisChangeTask(redisConnection redis.Conn, funds []string) *RedisChangeTask {
	return &RedisChangeTask{Base: &routingpool.Base{Name: "Redis Change Task", Response: make(chan bool)}, redis:redisConnection, funds:funds, encoder:redisLangEncoder}
}

func NewRedisPushTask() *RedisPushTask {
	return &RedisPushTask{message : make(chan interface{}, 1024), Base:&routingpool.Base{Name: "Redis Push Task", Response: make(chan bool)}, redis:RedisConnection, encoder:redisLangEncoder}
}

func (r *RedisManager) Run(id int) {
	r.caller(id)
}

func (r *RedisManager) caller(id int) {
	cntPush := 0
	cntChange := 0
	cntChangeRouting := 0
	exit := false

	for !exit {
		select {
		case <-r.pushDone:
			cntPush = cntPush + 1
			if cntPush == viper.GetInt("redis.pushtask.count") {
				funds, _ := getFunds()
				cntChangeRouting = len(funds)
				if cntChangeRouting%10 !=0 {
					cntChangeRouting = cntChangeRouting /10 +1
				} else {
					cntChangeRouting = cntChangeRouting/10
				}

				// Push the tasks to routing pool based on how many funds required.
				for index := 0; index < cntChangeRouting; index++ {
					fixedIndex := index
					// If it's the last routing, it needs to only get the last elements in array, otherice it will panic here.
					if index == (cntChangeRouting-1) {
						fixedIndex = index*10 + (len(funds)-index*10)
					}

					routingpool.PutTask(NewRedisChangeTask(r.redis, funds[index*10 : fixedIndex]))
				}
			}
		case <-r.changeDone:
			cntChange = cntChange + 1
			if cntChange == cntChangeRouting {
				logger.Debug("Received change done")
				exit = true
			}
		}
	}

	if r.redis != nil {
		logger.Debug("Redis closed.")
		r.redis.Close()
	}
}

func (c *RedisChangeTask) Run(id int) {
	c.caller(id)
}

// Plus the replicated records to one and delete them from Redis, then re-insert the result record into Redis
func (c *RedisChangeTask) caller(id int) {
	var	temp map[string]string

	for _, fund := range c.funds {
		changedRecords := make(map[string]map[string]string)

		logger.Debugf("Routing-%d, Trying to get records from Redis for %s.", id, fund)
		key := c.encoder.ConvertString(fund)
		values, err := redis.Values(c.redis.Do("LRANGE", key, 0, -1))
		if err != nil {
			logger.Errorf("redis lrange failed:", err)
		}

		for _, v1 := range values {
			json.Unmarshal(v1.([]byte), &temp)
			logger.Debugf("Routing-%d, Processing record - code:%s, recorddata:%s, holdcount::%s, holdvalue:%s", id, temp["code"], temp["recorddate"], temp["holdcount"], temp["holdvalue"])

			// If the record replicated, plus the records on fields
			if _, ok := changedRecords[temp["recorddate"]]; ok {
				logger.Debugf("Routing-%d, Found replicated record on date %s", id, temp["recorddate"])

				//record := make(map[string]string)
				changedRecords[temp["recorddate"]]["code"] = temp["code"]
				changedRecords[temp["recorddate"]]["recorddate"] = temp["recorddate"]
				changedRecords[temp["recorddate"]]["holdcount"] = fmt.Sprintf("%.4f", utility.String2Folat64(temp["holdcount"])+utility.String2Folat64(changedRecords[temp["recorddate"]]["holdcount"]))
				changedRecords[temp["recorddate"]]["holdvalue"] = fmt.Sprintf("%.4f", utility.String2Folat64(temp["holdvalue"])+utility.String2Folat64(changedRecords[temp["recorddate"]]["holdvalue"]))
				logger.Debugf("Routing-%d, re-caculate (%s,%s) record on date %s", id, changedRecords[temp["recorddate"]]["holdcount"], changedRecords[temp["recorddate"]]["holdvalue"], temp["recorddate"])
			} else {
				changedRecords[temp["recorddate"]] = make(map[string]string)

				changedRecords[temp["recorddate"]]["code"] = temp["code"]
				changedRecords[temp["recorddate"]]["recorddate"] = temp["code"]
				changedRecords[temp["recorddate"]]["holdcount"] = temp["holdcount"]
				changedRecords[temp["recorddate"]]["holdvalue"] = temp["holdvalue"]

				logger.Debugf("Routing-%d, (%s, %s) record on date %s is not replicated", id, changedRecords[temp["recorddate"]]["holdcount"], changedRecords[temp["recorddate"]]["holdvalue"], temp["recorddate"])
			}
		}

		// Delete the fund from Redis
		_, err = c.redis.Do("DEL", key)
		if err != nil {
			logger.Errorf("Routing-%d,  Delete fund %s failed, %s", id, fund, err)
			continue
		}
		logger.Debugf("Routing-%d, Delete fund %s records from Redis successfully", id, fund)

		// Insert the re-calculated records to Redis again
		for _, v := range changedRecords{
			json_value, _ := json.Marshal(v)
			_, err := c.redis.Do("LPUSH", key, json_value)
			if err != nil {
				logger.Errorf("Routing-%d, Insert the re-calculated records failed for fund %s, %s", id, key,  err)
				continue
			}
		}
	}

	ChangeTaskDone()
}


func PushDataIntoRedis(msg interface{}) {
	RedisPush.message <- msg
}

func PushTaskDone() {
	RedisDispatcher.pushDone <- true
}

func ChangeTaskDone() {
	RedisDispatcher.changeDone <- true
}
func (p *RedisPushTask) caller(id int) {
	timeout := time.NewTimer(time.Second * time.Duration(viper.GetInt("redis.pushtask.timer")))
	exit := false
	funds, _ := getFunds()
	//encoder := mahonia.NewEncoder("gbk")

	for !exit {
		select {
			case data := <-p.message:
				logger.Infof("Analysis-task %d, received data", id)
				tmp := data.([]*htmlparser.JJCCData)

				for _, value := range tmp {
					logger.Debugf("Analysis-task %d, row data name %s, code %s, holdcount %.4f, holdvalue %.4f", id, value.Name, value.Code, value.HoldCount, value.HoldValue)
					for _,fund := range funds {
						if strings.Contains(value.Name, fund) {
							logger.Debugf("Analysis-task %d, found record data for %s", id, fund)

							raw := map[string]string{	"code" : value.Code,
														"recorddate" : value.RecordDate,
														"holdcount" : fmt.Sprintf("%.4f", value.HoldCount),
														"holdvalue" : fmt.Sprintf("%.4f", value.HoldValue),
													}

							key := p.encoder.ConvertString(fund)
							json_value, _ := json.Marshal(raw)
							_, err := p.redis.Do("LPUSH", key, json_value)
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

func (p *RedisPushTask) Run(id int) {
	p.caller(id)
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