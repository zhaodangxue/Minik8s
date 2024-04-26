package apiobjects

import "time"

type TypeMeta struct {
	ApiVersion string
	Kind string
}

type ObjectMeta struct {
	Name string
	Namespace string
	Labels map[string]string
	UID string
	CreationTimestamp time.Time
	DeletionTimestamp time.Time
}
