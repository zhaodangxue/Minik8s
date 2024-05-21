package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/kubedns/nginx"

	log "github.com/sirupsen/logrus"
)


func getServiceAddr(serviceName string, namespace string) (string, error) {
	var service apiobjects.Service

	url :="/api/service/" + namespace + "/" + serviceName
	val, err := etcd.Get(url)
	if val == "" || err != nil{
		log.Error("[getServiceAddr] error getting service: ", err)
		return "", err
	}
	err = json.Unmarshal([]byte(val), &service)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if service.Data.Name == serviceName {
		//检查端口是否一致
		if service.Status.ClusterIP != "" {
			return service.Status.ClusterIP, nil
		}
	}

	return "", errors.New("[getServiceAddr] service not found")
}

func generatePath(rawPath string, host string, method string) error {
	parts := strings.Split(rawPath, ".")
	result := "/dns"
	for i := len(parts) - 1; i >= 0; i-- {
		result = result + "/" + parts[i]
	}
	log.Info("[generatePath] the new dns path is ", result)
	hoststr := apiobjects.DNSEntry{
		Host: "192.168.1.12",
	}

	if method == "create" {
		//err := dnsStorageTool.Create(context.Background(), result, &hoststr)
        jsonValue, err := json.Marshal(hoststr)
		if err != nil {
			return err
		}
		err = etcd.Put(result, string(jsonValue))
		return err
	} else {
		// 更新Dns
		return errors.New("[generatePath] method not supported")
	}
}

func updateNginx() error {
	allRecord := getAllDNSRecords()
	nginx.GenerateConfig(allRecord)
	return nil
	//err := nginx.ReloadNginx()
	//return err
}

func getAllDNSRecords() []apiobjects.DNSRecord {
	var dnsRecords []apiobjects.DNSRecord
	values, err := etcd.Get_prefix("/api/dns")
	if err != nil {
		log.Error("[getAllDNSRecords] error getting all DNS records: ", err)
	}
	for _, value := range values {
		var dnsRecord apiobjects.DNSRecord
		err := json.Unmarshal([]byte(value), &dnsRecord)
		if err != nil {
			fmt.Println(err)
			log.Error("[getAllDNSRecords] error getting all DNS records: ", err)
		}
		dnsRecords = append(dnsRecords, dnsRecord)
	}
	return dnsRecords
}