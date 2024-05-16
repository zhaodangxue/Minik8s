//go:build release
package node

import "time"

const (
	// 添加超时间隔的时间(ListCheckNode的调用间隔)
	NODE_TIMEOUT = 10 * time.Second
	// 超时次数上限
	NODE_TIMEOUT_COUNT_LIMIT = 3
	// 更新节点信息的时间间隔
	NODE_LIST_INTERVAL = 60 * time.Second
)