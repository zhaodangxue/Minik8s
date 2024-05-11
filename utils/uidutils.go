package utils

import uuid "github.com/satori/go.uuid"

func NewUUID() string {
	return uuid.NewV4().String()
}
