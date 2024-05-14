package api

import "minik8s/apiobjects"

type NodeHealthReportRequest struct {
	Node apiobjects.Node `json:"node"`
	Pods []*apiobjects.Pod `json:"pods"`
}

type NodeHealthReportResponse struct {
	// UnmatchedPods 用于存放未匹配的Pod的Path
	UnmatchedPodPaths []string `json:"unmatchedPodPaths"`
}
