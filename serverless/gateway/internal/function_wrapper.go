package internal

import (
	"minik8s/apiobjects"
	"sync/atomic"
)

type FunctionWrapper struct {
	Function *apiobjects.Function
	// 不使用大锁, 因为每个Function Call都会更改
	// 细粒度保证高并发
	QPSCounter *atomic.Int64
	// 缓存一份副本数，用于加速读
	// 由于Gateway是Function的ReplicaSet副本数量的唯一写者
	// 所以不用担心缓存一致性问题
	ScaleTarget int
}
