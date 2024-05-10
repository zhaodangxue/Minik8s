package ctlutils

import (
	"io"
	"os"
	"strings"
)

type ApiObjectType byte

const (
	Unknown ApiObjectType = iota //从0开始，依次加1
	Test
	Pod
	Node
	Pv
	Pvc
)

func (a ApiObjectType) String() string {
	switch a {
	case Test:
		return "test"
	case Pod:
		return "pod"
	case Node:
		return "node"
	case Pv:
		return "pv"
	case Pvc:
		return "pvc"
	default:
		return "unknown"
	}
}
func IsLetter(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}
func ParseApiObjectType(s []byte) ApiObjectType {
	sum_length := len(s)
	index := strings.Index(string(s), "kind:") + 5
	//跳过空格
	for index < sum_length && s[index] == ' ' {
		index++
	}
	start, end := index, index
	for end < sum_length && IsLetter(s[end]) {
		end++
	}
	tp := string(s[start:end])
	tp = strings.ToLower(tp)
	switch tp {
	case "test":
		return Test
	case "pod":
		return Pod
	case "node":
		return Node
	case "persistentvolume":
		return Pv
	case "persistentvolumeclaim":
		return Pvc
	default:
		return Unknown
	}
}
func LoadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func ParseType(path string) (ApiObjectType, error) {
	data, err := LoadFile(path)
	if err != nil {
		return Unknown, err
	}
	return ParseApiObjectType(data), nil
}
