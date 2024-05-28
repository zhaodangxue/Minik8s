package serveless_utils

import "minik8s/apiobjects"

type NodeType byte

const (
	BranchType NodeType = iota
	FunctionType
)

type node struct {
	Type NodeType
	Name string
	*function
	*choices
}
type function struct {
	Next *node
}
type choices struct {
	Branches []*branch
}
type branch struct {
	Variable   string
	Next       *node
	BranchFunc branchFunc
}
type branchFunc func(data string, variable string, result interface{}) bool
type DAG struct {
	Root *node
}

func (c *choices) chooseBranch(data string, result interface{}) *node {
	for _, br := range c.Branches {
		if br.BranchFunc(data, br.Variable, result) {
			return br.Next
		}
	}
	return nil
}
func Workflow2DAG(wf *apiobjects.Workflow) *DAG {
	nodeMap := make(map[string]apiobjects.WorkflowNode)
	dagMap := make(map[string]*node)
	for name := range wf.Nodes {
		nodeMap[name] = wf.Nodes[name]
		dagMap[name] = &node{Name: name}
	}
	if _, exist := dagMap[wf.Begin]; !exist {
		return nil
	}
	root := BuildDAG(wf.Begin, dagMap, nodeMap)
	if root == nil {
		return nil
	}
	return &DAG{
		Root: root,
	}
}
func chooseJudgeFunction(b apiobjects.Branch) branchFunc {
	switch {
	case b.IntegerEqual != nil:
		return IntegerEqual
	case b.IntegerNotEqual != nil:
		return IntegerNotEqual
	case b.IntegerGreaterThan != nil:
		return IntegerGreaterThan
	case b.IntegerLessThan != nil:
		return IntegerLessThan
	case b.BooleanEqual != nil:
		return BooleanEqual
	case b.BooleanNotEqual != nil:
		return BooleanNotEqual
	case b.StringEqual != nil:
		return StringEqual
	case b.StringNotEqual != nil:
		return StringNotEqual
	case b.FloatEqual != nil:
		return FloatEqual
	case b.FloatNotEqual != nil:
		return FloatNotEqual
	case b.FloatGreaterThan != nil:
		return FloatGreaterThan
	case b.FloatLessThan != nil:
		return FloatLessThan
	}
	return nil
}
func BuildDAG(currentNode string, dagMap map[string]*node, nodeMap map[string]apiobjects.WorkflowNode) *node {
	wfNode := nodeMap[currentNode]
	dagNode := dagMap[currentNode]
	if dagNode == nil {
		return nil
	}
	var next *node
	switch wfNode.Type {
	case apiobjects.NodeTypeFunction:
		if wfNode.FunctionNode != nil && wfNode.Next != nil {
			next = BuildDAG(*wfNode.Next, dagMap, nodeMap)
		} else {
			next = nil
		}
		return &node{
			Type: FunctionType,
			Name: currentNode,
			function: &function{
				Next: next,
			},
			choices: nil,
		}
	case apiobjects.NodeTypeBranch:
		branchs := wfNode.Branchs
		var branches []*branch
		if branchs != nil {
			for _, b := range branchs.Branchs {
				if b.Next != nil {
					next = BuildDAG(*b.Next, dagMap, nodeMap)
				} else {
					next = nil
				}
				br := &branch{
					Variable:   b.Variable,
					Next:       next,
					BranchFunc: chooseJudgeFunction(b),
				}
				branches = append(branches, br)
			}
		}
		return &node{
			Type:     BranchType,
			Name:     currentNode,
			function: nil,
			choices: &choices{
				Branches: branches,
			},
		}
	}
	return nil
}
