package redisstorage

import (
	runt "runtime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/redis/v3"
	"xn--gckvb8fzb.com/glides/services/config"
	"xn--gckvb8fzb.com/hyperuplink/errs"
)

func New(cfg *config.Config) (storage fiber.Storage, err error) {
	rcfg, err := cfg.Redis()
	if err != nil {
		return nil, err
	}

	redisConfig := redis.Config{
		MasterName: rcfg.MasterName,
		Username:   rcfg.Username,
		Password:   rcfg.Password,
		Database:   rcfg.Database,
		Reset:      rcfg.Reset,
	}

	switch len(rcfg.Addrs) {
	case 0:
		return nil, errs.ErrRedisAddrsEmpty
	case 1:
		hostPort := strings.Split(rcfg.Addrs[0], ":")
		if len(hostPort) != 2 {
			return nil, errs.ErrRedisAddrsMalformed
		}
		redisConfig.Host = hostPort[0]
		if redisConfig.Port, err = strconv.Atoi(hostPort[1]); err != nil {
			return nil, errs.ErrRedisAddrsMalformed
		}
	default:
		redisConfig.Addrs = rcfg.Addrs
	}

	poolSize := rcfg.Poolsize
	if poolSize <= 0 {
		poolSize = 10 * runt.GOMAXPROCS(0)
	}
	redisConfig.PoolSize = poolSize

	return redis.New(redisConfig), nil
}
