// Description: This file contains the definition of the Object, ObjectMeta, ObjectRef and TypeMeta structs.
// These structs are used to define the basic structure of the objects that are used in the Kubernetes API.
package apiobjects

import "time"

type TypeMeta struct {
	ApiVersion string
	Kind string
}

type ObjectMeta struct {
	Name string
	Namespace string
	UID string
	Labels map[string]string
	CreationTimestamp time.Time
	DeletionTimestamp time.Time
}

type Object struct {
	TypeMeta
	ObjectMeta
}

func (obj *Object) GetObjectRef() ObjectRef {
	return ObjectRef{
		TypeMeta: obj.TypeMeta,
		Name: obj.Name,
		Namespace: obj.Namespace,
		UID: obj.UID,
	}
}

func (obj *Object) GetObjectPath() string {
	return obj.ApiVersion + "/" + obj.Kind + "/" + obj.Namespace + "/" + obj.Name
}

type ObjectRef struct {
	TypeMeta
	Name      string
	Namespace string
	UID       string
}

func (ref *ObjectRef) GetObjectPath() string {
	return ref.ApiVersion + "/" + ref.Kind + "/" + ref.Namespace + "/" + ref.Name
}
