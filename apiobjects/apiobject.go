// Description: 定义了Kubernetes API对象的基本结构。
package apiobjects

import (
	"strings"
	"time"
)

type ActionType byte

// TODO: give kind a seperate type

const (
	Create ActionType = iota //从0开始，依次加1
	Update
	Delete
)

type TopicMessage struct {
	ActionType ActionType
	Object     string
}

type TypeMeta struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type ObjectMeta struct {
	Name              string            `yaml:"name"`
	Namespace         string            `yaml:"namespace"`
	UID               string            `yaml:"uid"`
	Labels            map[string]string `yaml:"labels"`
	CreationTimestamp time.Time         `yaml:"-"`
	DeletionTimestamp time.Time         `yaml:"-"`
}

// Object is the base struct for all objects in the Kubernetes API.
// 包含TypeMeta和ObjectMeta
// 可以使用GetObjectRef从Object中获取ObjectRef。
// 可以使用GetObjectPath获取Object的路径。
type Object struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata"`
}

func (obj *Object) GetObjectRef() ObjectRef {
	return ObjectRef{
		TypeMeta:  obj.TypeMeta,
		Name:      obj.Name,
		Namespace: obj.Namespace,
		UID:       obj.UID,
	}
}
func (obj *Object) GetObjectPath() string {
	return "/api" + "/" + strings.ToLower(obj.Kind) + "/" + obj.Namespace + "/" + obj.Name
}

// 可以唯一标识一个对象的引用。
// 可以使用GetObjectPath获取对象的路径。
type ObjectRef struct {
	TypeMeta
	Name      string
	Namespace string
	UID       string
}

func (ref *ObjectRef) GetObjectPath() string {
	return ref.ApiVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name
}

type Base_test struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}
type TestYaml struct {
	Base_test `yaml:",inline"`
	Spec      test_spec `yaml:"spec"`
}
type test_spec struct {
	Replicas int32  `yaml:"replicas"`
	Name     string `yaml:"name"`
}
