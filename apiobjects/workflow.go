package apiobjects

type NodeType string

const (
	NodeTypeFunction NodeType = "function"
	NodeTypeBranch   NodeType = "branch"
)

type Workflow struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Begin      string                  `json:"begin"`
	Nodes      map[string]WorkflowNode `json:"nodes"`
}

type WorkflowNode struct {
	Type          NodeType `json:"type"`
	*FunctionNode `json:",inline"`
	*Branchs      `json:",inline"`
}

type FunctionNode struct {
	Next *string `json:"next"`
}

type Branchs struct {
	Branchs []Branch `json:"branchs"`
}
type Branch struct {
	Variable           string   `json:"variable"`
	Next               *string  `json:"next"`
	IntegerEqual       *int64   `json:"integerEqual,omitempty"`
	IntegerNotEqual    *int64   `json:"integerNotEqual,omitempty"`
	IntegerLessThan    *int64   `json:"integerLessThan,omitempty"`
	IntegerGreaterThan *int64   `json:"integerGreaterThan,omitempty"`
	BooleanEqual       *bool    `json:"booleanEqual,omitempty"`
	BooleanNotEqual    *bool    `json:"booleanNotEqual,omitempty"`
	StringEqual        *string  `json:"stringEqual,omitempty"`
	StringNotEqual     *string  `json:"stringNotEqual,omitempty"`
	FloatEqual         *float64 `json:"floatEqual,omitempty"`
	FloatNotEqual      *float64 `json:"floatNotEqual,omitempty"`
	FloatLessThan      *float64 `json:"floatLessThan,omitempty"`
	FloatGreaterThan   *float64 `json:"floatGreaterThan,omitempty"`
}
