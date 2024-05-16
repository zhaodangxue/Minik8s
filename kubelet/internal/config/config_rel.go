//go:build release

package config

import "time"

// HealthReportInterval 用于定时上报Node和Pod的状态
const (
	HealthReportInterval = 10 * time.Second
)

