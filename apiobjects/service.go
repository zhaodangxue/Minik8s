package apiobjects

import "encoding/json"

// a sevice example for k8s service
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
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	UID       string `json:"uid,omitempty" yaml:"uid,omitempty"` // 一个service的唯一标识
}

type Service struct {
	APIVersion string   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"` // api版本
	Data       MetaData `json:"metadata" yaml:"metadata"`

	// 定义了service的规范
	Spec ServiceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	// 表示service的状态
	Status ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}
type ServiceSpec struct {
	// 只有 ClusterIP 类型
	Type ServiceType `json:"type,omitempty" yaml:"type,omitempty"`

	// 这个服务的所有端口
	Ports []ServicePort `json:"ports" yaml:"ports"`

	//通过这个selector来选择pod
	Selector map[string]string `json:"selector" yaml:"selector"`
}

// service的一个端口
type ServicePort struct {
	//一个service可以有多个端口，每个端口都有一个不同的name
	Name string `json:"name" yaml:"name"`

	// 协议
	Protocol Protocol `json:"protocol" yaml:"protocol"`

	// 对外暴露的端口
	Port int32 `json:"port" yaml:"port"`

	// 转发的端口，pod对应的端口
	TargetPort string `json:"targetPort" yaml:"targetPort"`
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
	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`

	//ClusterIP是service的IP地址，通常是由master随机分配的
	ClusterIP string `json:"clusterIP" yaml:"clusterIP"`
}

func (s *Service) GetType() string {
	return "service"
}

func (s *Service) GetObjectPath() string {
	return "/api/service/" + s.Data.Namespace + "/" + s.Data.Name
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
