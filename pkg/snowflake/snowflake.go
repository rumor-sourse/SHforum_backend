package snowflake

import (
	"SHforum_backend/interval/settings"
	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"log"
	"time"
)

func Init(startTime string) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	return
}

func GenID(tableName string) int64 {
	config, ok := settings.Conf.SnowFlakeConfig.TableConfig[tableName]
	if !ok {
		zap.L().Error("tableName not found")
	}
	node, err := sf.NewNode(int64(config.DataCenterID<<5 | config.WorkerID))
	if err != nil {
		log.Fatalf("Failed to create snowflake node: %v", err)
	}
	return node.Generate().Int64()
}
