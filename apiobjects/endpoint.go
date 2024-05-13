package apiobjects

import (
	"encoding/json"
)

type Endpoint struct {
	ServiceName string       `json:"serviceName" yaml:"serviceName"`
	Data        MetaData     `json:"metadata" yaml:"metadata"`
	Spec        EndpointSpec `json:"spec" yaml:"spec"`
}

type EndpointSpec struct {
	SvcIP    string `json:"svcIP" yaml:"svcIP"`
	SvcPort  int32  `json:"svcPort" yaml:"svcPort"`
	DestIP   string `json:"dstIP" yaml:"dstIP"`
	DestPort int32  `json:"dstPort" yaml:"dstPort"`
}

func (e *Endpoint) GetType() string {
	return "endpoint"
}

func (e *Endpoint) GetObjectPath() string {
	return "/api/endpoint/" + e.ServiceName + "/" + e.Data.Namespace + "/" + e.Data.Name
}

func (e *Endpoint) MarshalJSON() ([]byte, error) {
	type Alias Endpoint
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

func (e *Endpoint) UnMarshalJSON(data []byte) error {
	type Alias Endpoint
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) Union(other *Endpoint) {

}
