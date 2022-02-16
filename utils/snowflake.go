package utils

import (
	"time"

	"github.com/bwmarrin/snowflake"
	log "github.com/truxcoder/truxlog"
)

var node *snowflake.Node

func init() {
	var nodeId int64
	nodeId = 1
	st, err := time.Parse("2006-01-02", "2021-11-11")
	if err != nil {
		log.Error(err)
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	snowflake.NodeBits = 2
	snowflake.StepBits = 8
	node, err = snowflake.NewNode(nodeId)
	if err != nil {
		log.Error(err)
		return
	}
}

func GenId() int64 {
	id := node.Generate()
	return id.Int64()
}
