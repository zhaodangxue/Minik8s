package apiobjects

import "encoding/json"

// a sevice struct for k8s service
/*
apiVersion: v1
kind: Service
metadata:
  name: service-practice
  namespace: default
spec:
  selector:
    app: deploy-practice
  type: ClusterIP
  ports:
  - name: service-port1
    protocol: TCP
    port: 8080 # 对外暴露的端口
    targetPort: p1 # 转发的端口，pod对应的端口
  - name: service-port2
    protocol: TCP
    port: 3000 # 对外暴露的端口
    targetPort: p2 # 转发的端口，pod对应的端口
*/
type MetaData struct {
	Name             string            `json:"name,omitempty"`
	Namespace        string            `json:"namespace,omitempty"`
	UID              string            `yaml:"uid"`
	Labels           map[string]string `json:"labels,omitempty"`
}

type Service struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind"`
	Data MetaData `json:"metadata"`

	// 定义了service的规范
	Spec ServiceSpec `json:"spec,omitempty"`

	// 表示service的状态
	Status ServiceStatus `json:"status,omitempty"`
}
type ServiceSpec struct {
	// 只有 ClusterIP 类型
	Type ServiceType `json:"type,omitempty"`

	// 这个服务的所有端口
	Ports []ServicePort `json:"ports"`

	//通过这个selector来选择pod
	Selector map[string]string `json:"selector"`
}

// service的一个端口
type ServicePort struct {
	//一个service可以有多个端口，每个端口都有一个不同的name
	Name string `json:"name"`

	// 代表了service的IP协议
	Protocol Protocol `json:"protocol"`

	// 对外暴露的端口
	Port int32 `json:"port"`

    // 转发的端口，pod对应的端口
	TargetPort string `json:"targetPort"`
}

type ServiceType string
const (
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	//这个类型的service会被暴露到每个node上，是ClusterIP的扩展
	ServiceTypeNodePort ServiceType = "NodePort"

	//这个类型的service会被暴露到外部的负载均衡器上，是NodePort的扩展
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
)

type Protocol string
const (
	ProtocolTCP Protocol = "TCP"

	ProtocolUDP Protocol = "UDP"

	ProtocolSCTP Protocol = "SCTP"
)

type ServiceStatus struct {
	/*
		CREATING: 等待分配cluster ip
		CREATED: cluster ip分配完成
	*/
	Phase string `json:"phase,omitempty"`

	//ClusterIP是service的IP地址，通常是由master随机分配的
	ClusterIP string `json:"clusterIP"`
}

func (s *Service) UnMarshalJSON(data []byte) error {
	type Alias Service
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (s *Service) MarshalJSON() ([]byte, error) {
	type Alias Service
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

func (s *Service) Union(other *Service) {
	if s.Status.Phase == "" {
		s.Status.Phase = other.Status.Phase
	}
	if s.Status.ClusterIP == "" {
		s.Status.ClusterIP = other.Status.ClusterIP
	}
}