package node

import (
	"encoding/json"
	"errors"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/controller/api"
	"minik8s/global"
	"minik8s/utils"
	"sync"
	"time"
)

/*
管理节点的控制器

在内存中保持一个PodPath ->（Pod，超时计数）的map，用于存储所有的节点信息

有两个List函数

1. ListNode 用于定期从etcd获取节点信息，更新内存中的map

2. ListCheckNode 用于给Node添加超时计数，如果超时计数超过一定数量则删除

有一个Watch函数

 1. WatchNode 用于监听Node HealthReport消息，如果有节点的HealthReport消息，则清除对应节点的超时计数
    如果有新的节点加入，则添加到内存中的map中

由于操作间隔很大，所以用一把锁map的大锁保证线程安全
*/
type NodeController struct {
	InitInfo          api.InitStruct
	NodeList          map[string]nodeWrapper
	NodeListLock      sync.Mutex
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
}

const (
	// 添加超时间隔的时间(ListCheckNode的调用间隔)
	NODE_TIMEOUT = 10 * time.Second
	// 超时次数上限
	NODE_TIMEOUT_COUNT_LIMIT = 3
	// 更新节点信息的时间间隔
	NODE_LIST_INTERVAL = 60 * time.Second
)

type nodeWrapper struct {
	Node         *apiobjects.Node
	TimeoutCount int
}

func (c *NodeController) Init(init api.InitStruct) {
	c.InitInfo = init

	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     ListNode,
		Interval: NODE_LIST_INTERVAL,
	})
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     ListCheckNode,
		Interval: NODE_TIMEOUT,
	})

	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  WatchNode,
		Topic: global.NodeStateTopic(),
	})

	c.NodeList = make(map[string]nodeWrapper)
}

func (c *NodeController) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}

func (c *NodeController) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}

func WatchNode(c api.Controller, message apiobjects.TopicMessage) (err error) {
	utils.Debug("NodeController: WatchNode: ", message)
	node := new(apiobjects.Node)
	err = json.Unmarshal([]byte(message.Object), node)
	if err != nil {
		return
	}

	nodeController := c.(*NodeController)

	if message.ActionType != apiobjects.Update {
		return errors.New("NodeController: WatchNode: ActionType is not Update")
	}

	nodeController.NodeListLock.Lock()
	defer nodeController.NodeListLock.Unlock()
	if _, ok := nodeController.NodeList[node.GetObjectPath()]; !ok {
		// 如果是新节点，则添加到内存中的map中
		nodeController.NodeList[node.GetObjectPath()] = nodeWrapper{
			Node:         node,
			TimeoutCount: 0,
		}
		utils.Info("NodeController: WatchNode: New node added: ", node.GetObjectPath())
	} else {
		// 如果是已有节点，则清除超时计数
		nodeController.NodeList[node.GetObjectPath()] = nodeWrapper{
			Node:         node,
			TimeoutCount: 0,
		}
		utils.Info("NodeController: WatchNode: Node health report: ", node.GetObjectPath())
	}

	return
}

// ListNode 用于定期从etcd获取节点信息，更新内存中的map
func ListNode(c api.Controller) (err error) {
	utils.Debug("NodeController: ListNode")
	nodeController := c.(*NodeController)

	apiserverUrl := nodeController.InitInfo.GetApiserverUrl()
	var nodeList []*apiobjects.Node = make([]*apiobjects.Node, 0)
	err = utils.GetUnmarshal(apiserverUrl+route.NodePath, &nodeList)
	if err != nil {
		return
	}

	nodeController.NodeListLock.Lock()
	defer nodeController.NodeListLock.Unlock()
	for _, node := range nodeList {
		if _, ok := nodeController.NodeList[node.GetObjectPath()]; !ok {
			nodeController.NodeList[node.GetObjectPath()] = nodeWrapper{
				Node:         node,
				TimeoutCount: 0,
			}
		} else {
			nodeController.NodeList[node.GetObjectPath()] = nodeWrapper{
				Node:         node,
				TimeoutCount: nodeController.NodeList[node.GetObjectPath()].TimeoutCount,
			}
		}
	}
	utils.Info("NodeController: ListNode: Node list updated. Current size: ", len(nodeList))

	return
}

// ListCheckNode 用于给Node添加超时计数，如果超时计数超过一定数量则删除
func ListCheckNode(c api.Controller) (err error) {
	utils.Debug("NodeController: ListCheckNode")
	nodeController := c.(*NodeController)

	nodeController.NodeListLock.Lock()
	defer nodeController.NodeListLock.Unlock()
	for key, nw := range nodeController.NodeList {
		if nw.TimeoutCount >= NODE_TIMEOUT_COUNT_LIMIT {
			utils.Info("NodeController: ListCheckNode: Node timeout: ", nw.Node.GetObjectPath())
			// 在服务端删除
			_, err := utils.Delete(nodeController.InitInfo.GetApiserverUrl() + nw.Node.GetObjectPath())
			// 在本地删除
			// CHECK: 这个删除和检查error的顺序会否带来问题？
			delete(nodeController.NodeList, key)
			if err != nil {
				utils.Error("NodeController: ListCheckNode: Delete error: ", err)
				continue
			}
			utils.Info("NodeController: ListCheckNode: Node deleted: ", nw.Node.GetObjectPath())
		} else {
			// 超时计数加一
			nodeController.NodeList[key] = nodeWrapper{
				Node:         nw.Node,
				TimeoutCount: nw.TimeoutCount + 1,
			}
		}
	}

	return
}
