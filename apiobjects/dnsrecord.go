package apiobjects

// example:

import "encoding/json"

type DNSRecord struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"` 
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Name       string `json:"name" yaml:"name"`
	NameSpace 	string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Host       string `json:"host" yaml:"host"`
	Paths      []Path `json:"paths" yaml:"paths"`
}

type Path struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	PathName string `json:"pathName,omitempty" yaml:"pathName,omitempty"`
	Service string `json:"service" yaml:"service"`
	Port    int    `json:"port" yaml:"port"`
}

type DNSEntry struct {
	Host string `json:"host" yaml:"host"`
}

func (r *DNSRecord) GetObjectPath() string {
	return "/api/dns/"+ r.NameSpace + "/" + r.Name
}

func (r *DNSRecord) MarshalJSON() ([]byte, error) {
	type Alias DNSRecord
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

func (r *DNSRecord) UnMarshalJSON(data []byte) error {
	type Alias DNSRecord
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
