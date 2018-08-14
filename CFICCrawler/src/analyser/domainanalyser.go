package analyser

import (
	"github.com/garyburd/redigo/redis"
	"crawler/dataminer/database"
	"utility"
)

var logger = utility.GetLogger()

func DoDomainAnalyse() {

	domains, err := redis.Strings(database.RedisPool.Get().Do("smembers", "DOMAINS"))
	if len(domains) == 0 || err != nil {
		logger.Errorf("Query Redis failed for DOMAINS")
		return
	}

}
