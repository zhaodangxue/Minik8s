package etcd

import (
	"context"
	"log"
	"minik8s/global"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *clientv3.Client
	config clientv3.Config
	ctx    = context.Background()
	Err    error
)

func init() {
	config = clientv3.Config{
		Endpoints:   []string{global.EtcdAndRedisHost + ":2379"},
		DialTimeout: 5 * time.Second,
	}
	client, Err = clientv3.New(config)
	if Err != nil {
		log.Fatal(Err)
	}
}
func CheckClient() (err error) {
	if client == nil {
		client, err = clientv3.New(config)
		if err != nil {
			return err
		}
	}
	return nil
}
func Put(key, value string) (err error) {
	if err = CheckClient(); err != nil {
		log.Fatal(err)
		return err
	}
	_, err = client.Put(ctx, key, value)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func Get(key string) (value string, err error) {
	if err = CheckClient(); err != nil {
		log.Fatal(err)
		return "", err
	}
	resp, err := client.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}
func Delete(key string) (err error) {
	if err = CheckClient(); err != nil {
		log.Fatal(err)
		return err
	}
	_, err = client.Delete(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func Delete_prefix(prefix string) (err error) {
	if err = CheckClient(); err != nil {
		log.Fatal(err)
		return err
	}
	_, err = client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func Get_prefix(prefix string) (values []string, err error) {
	if err = CheckClient(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, ev := range resp.Kvs {
		values = append(values, string(ev.Value))
	}
	return values, nil
}
func Clear() error {
	return Delete_prefix("")
}
