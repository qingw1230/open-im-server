package getcdv3

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type RegEtcd struct {
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

var rEtcd *RegEtcd

// "schema:///serviceName/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "schema:///serviceName"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl)
}

func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	cli, err := clientv3.New(clientv3.Config{Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})
	if err != nil {
		return fmt.Errorf("create etcd clientv3 client failed, errmsg: %v, etcd addr: %s", err, etcdAddr)
	}
	fmt.Println("RegisterEtcd")

	// lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		return fmt.Errorf("grant failed")
	}

	// schema:///serviceName/ip:port -> ip:port
	serviceValue := net.JoinHostPort(myHost, strconv.Itoa(myPort))
	serviceKey := GetPrefix(schema, serviceName) + serviceValue
	fmt.Printf("serviceKey: %s, serviceValue: %s\n", serviceKey, serviceValue)
	// set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("put failed, errmsg:%v, key:%s, value:%s", err, serviceKey, serviceValue)
	}

	// keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		return fmt.Errorf("keepalive failed, errmsg:%v, lease id:%d", err, resp.ID)
	}
	fmt.Println("RegisterEtcd ok")
	go func() {
		for {
			select {
			case v, ok := <-kresp:
				if ok == true {
					// fmt.Println(" kresp ok ", v)
				} else {
					fmt.Println(" kresp failed ", v)
				}
			}
		}
	}()

	rEtcd = &RegEtcd{ctx: ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey,
	}
	return nil
}

func UnRegisterEtcd() {
	// delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}
