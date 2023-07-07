package snowflake

import (
	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/snowflake"
)

var snowFlakeWorker *snowflake.Node

func init() {
	logger.InfoF("new snowflake worker id = %d", conf.Cfg.Snowflake.WorkerId)

	worker, err := snowflake.NewNode(
		conf.Cfg.Snowflake.WorkerId,
	)

	if err != nil {
		panic(err)
	}

	snowFlakeWorker = worker
}

func Id() int64 {
	if snowFlakeWorker != nil {
		return int64(snowFlakeWorker.Generate())
	}

	return 0
}

func IdBatch(length int) []int64 {
	resId := make([]int64, 0)
	if snowFlakeWorker != nil {
		i := 0
		for {
			if i >= length {
				break
			}
			resId = append(resId, int64(snowFlakeWorker.Generate()))
			i++
		}
	}

	return resId
}
