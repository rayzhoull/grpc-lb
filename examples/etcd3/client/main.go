package main

import (
	etcd "github.com/coreos/etcd/clientv3"
	grpclb "github.com/rayzhoull/grpc-lb"
	registry "github.com/rayzhoull/grpc-lb/registry/etcd3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
	"github.com/rayzhoull/grpc-lb/examples/proto"
)

func main() {
	etcdConfg := etcd.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	r := registry.NewResolver("/grpc-lb", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewRoundRobinSelector())
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithTimeout(time.Second*5))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()

	client := proto.NewTestClient(c)

	for i := 0; i < 50; i++ {
		log.Println("i:",i)
		resp, err := client.Say(context.Background(), &proto.SayReq{Content: "round robin"})
		if err != nil {
			log.Println("aa:", err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		log.Printf(resp.Content)
	}

}
