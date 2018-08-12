package analyzer

import (
	"routingpool"
	"htmlparser"
	"github.com/spf13/viper"
	"time"
	"os"
	"bufio"
	"github.com/axgle/mahonia"
	"fmt"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"utility"
	"flag"
)

var (
	RedisPusher *RedisPusherTask
	RedisDispatcher *RedisManager

	redisServer = flag.String("127.0.0.1", ":6379", "")
	redisPool *redis.Pool
	RedisConnection redis.Conn
	redisLangEncoder mahonia.Encoder
)

func init() {
	redisPool = newPool(*redisServer)
	RedisConnection = initRedis()
	redisLangEncoder = mahonia.NewEncoder("gbk")

	RedisPusher = NewRedisPushTask()
	RedisDispatcher = NewRedisManagerTask()
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func initRedis() redis.Conn {
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
	stockname chan []string  // format example: [SHH, stockname]
	domains chan map[string][]string // format example: [601111]{d1, d2, d3, d4, d5}
	redis redis.Conn
	pool *redis.Pool
}

// Implement Task interface
type RedisPusherTask struct {
	*routingpool.Base
	jjcc chan interface{}

	pool *redis.Pool
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
	return &RedisManager{pushDone:make(chan bool),
						  changeDone:make(chan bool),
						  stockname:make(chan []string),
						  domains:make(chan map[string][]string),
						  redis:RedisConnection,
						  pool : redisPool,
						  Base: &routingpool.Base{Name: "Redis Change Task", Response: make(chan bool)}}
}

func NewRedisChangeTask(redisConnection redis.Conn, funds []string) *RedisChangeTask {
	return &RedisChangeTask{Base: &routingpool.Base{Name: "Redis Change Task", Response: make(chan bool)}, redis:redisConnection, funds:funds, encoder:redisLangEncoder}
}

func NewRedisPushTask() *RedisPusherTask {
	return &RedisPusherTask{jjcc : make(chan interface{}, 1024),
							 Base:&routingpool.Base{Name: "Redis Push Task", Response: make(chan bool)},
							 redis:RedisConnection,
							 encoder:redisLangEncoder,
							 pool:redisPool}
}

func PushStocks(name []string) {
	RedisDispatcher.stockname <- name
}

func PushDomains(domains map[string][]string)  {
	RedisDispatcher.domains <- domains
}

func (r *RedisManager) Run(id int) {
	r.caller(id)
}

func (r *RedisManager) caller(id int) {
	for pusher := 0; pusher < viper.GetInt("redis.pushtask.count"); pusher++ {
		routingpool.PutTask(RedisPusher)
	}

	cntPush := 0
	cntChangeRouting := 0
	exit := false

	for !exit {
		select {
		case <-r.pushDone:
			cntPush = cntPush + 1
			if cntPush == viper.GetInt("redis.pushtask.count") {
				exit = true
				break

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
		}
	}

	if redisPool != nil {
		logger.Debug("Redis-Pool closed.")
		redisPool.Close()
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
				changedRecords[temp["recorddate"]]["recorddate"] = temp["recorddate"]
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
		for _, v := range changedRecords {
			json_value, _ := json.Marshal(v)
			_, err := c.redis.Do("LPUSH", key, json_value)
			if err != nil {
				logger.Errorf("Routing-%d, Insert the re-calculated records failed for fund %s, %s", id, key,  err)
				continue
			}

			filename := fmt.Sprintf("d:/out/%s.csv", fund)
			line := fmt.Sprintf("%s,%s,%s", v["recorddate"], v["holdcount"], v["holdvalue"])
			utility.WriteToFile(filename, line)
		}
	}

	ChangeTaskDone()
}


func PushDataIntoRedis(msg interface{}) {
	RedisPusher.jjcc <- msg
}

func PushTaskDone() {
	RedisDispatcher.pushDone <- true
}

func ChangeTaskDone() {
	RedisDispatcher.changeDone <- true
}

func (p *RedisPusherTask) Run(id int) {
	p.caller(id)
}

func (p *RedisPusherTask) caller(id int) {
	timeout := time.NewTimer(time.Second * time.Duration(viper.GetInt("redis.pushtask.timer")))
	exit := false

	for !exit {
		select {
			case data := <-p.jjcc:
				records := data.([]*htmlparser.JJCCData)
				if records == nil {
					logger.Warning("Received empty message when processing JJCC data.")
					continue
				}

				temp2 := map[string]map[string]string{}
				for _, value := range records {
					logger.Debugf("RedisPusherTask %d, constructing raw data to json format for %s on %s:raw value is name %s, code %s, holdcount %.4f, holdvalue %.4f",
											id,
											value.Stock_name,
											value.RecordDate,
											value.Name,
											value.Code,
											value.HoldCount,
											value.HoldValue)

					// Add fund name and func code mapping to table 'FUND_INFO_TABLE' which is a set and keep the data uniqueness
					fund_info := p.encoder.ConvertString(records[0].Code + "_" + records[0].Name)
					_, err := p.pool.Get().Do("SADD", "FUND_INFO_TABLE", fund_info)
					if err != nil {
						logger.Errorf("Push fund info to Redis failure:", err)
					}

					temp1 := map[string]string{"count": fmt.Sprintf("%.4f", value.HoldCount), "value": fmt.Sprintf("%.4f", value.HoldValue)}
					temp2[value.Code] = temp1
				}

				rowData := map[string]map[string]map[string]string{records[0].RecordDate:temp2}

				/* JJCC data presentation in Redis
				{
				"2017-09-30":{"003594":{"count":"23.6700","value":"0.0138"},
							 "003641":{"count":"23.6700","value":"0.0138"},
							 "003922":{"count":"23.6700","value":"0.0138"},
							 "003924":{"count":"23.6700","value":"0.0138"},
							 "004336":{"count":"23.6700","value":"0.0138"},
							 "004338":{"count":"23.6700","value":"0.0138"},
							 "004434":{"count":"1210.0069","value":"0.7042"},
							 "050001":{"count":"3800.0000","value":"2.2116"},
							 "050023":{"count":"116.7300","value":"0.0679"},
							 "050201":{"count":"1450.0000","value":"0.8439"},
							 "160512":{"count":"1350.0012","value":"0.7857"}}}
				{
				"2017-12-31":{"160512":{"count":"1000.0012","value":"0.5210"}}}
				{
				"2018-03-31":{"004194":{"count":"27.0900","value":"0.0115"},
							  "160512":{"count":"800.0012","value":"0.3384"}}}
				*/
				key := p.encoder.ConvertString(records[0].Stock_name + "_" + records[0].Stock_number)
				json_value, _ := json.Marshal(rowData)
				_, err := p.pool.Get().Do("LPUSH", key, json_value)
				if err != nil {
					logger.Errorf("Push JJCC data to Redis failure:", err)
				}


				timeout.Reset(time.Second * time.Duration(viper.GetInt("analyser.timer")))
			case stock := <- RedisDispatcher.stockname:
				/* Stocks name presentation in Redis
				SHH:{name1_code, name2_code, name3_code, name4_code......}
				SHZ:(name1_code, name2_code, name3_code, name4_code......)
				CYB:{name1_code, name2_code, name3_code, name4_code......}
				 */

				// Encoding the value because the value contains Chinese so that check it in Redis directly.
				encoded_value := p.encoder.ConvertString(stock[1])
				_, err := p.pool.Get().Do("LPUSH", stock[0], encoded_value)
				if err != nil {
					logger.Errorf("Push stock name to Redis failure:", err)
				}

			case domains := <- RedisDispatcher.domains:
				for k, v := range domains  {
					json_value, _ := json.Marshal(v)

					_, err := p.pool.Get().Do("SET", "DOMAIN_" + k, json_value)
					if err != nil {
						logger.Errorf("Push Domain to Redis failure:", err)
					}
				}


			case <- timeout.C: // The waiting seconds for receiving data, analysis routine exits after the waiting.
				PushTaskDone()
				exit = true
				break
		}
	}

	logger.Infof("RedisPusherTask %d, exit.....................", id)
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