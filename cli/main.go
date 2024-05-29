package main

import (
	"context"
	"fmt"
	"log"
	"obzev0/shared/services/tcpService"
	"obzev0/shared/structs"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

func LoadConfig(filename string) (structs.Config, error) {
	var config structs.Config

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read YAML file: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}

func main() {
	conn, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := tcpService.NewTcpServiceClient(conn)

	cnf, err := LoadConfig("obzevConf.yaml")
	config := &tcpService.TcpConfig{
		ReqDelay: cnf.Delays.ReqDelay,
		ResDelay: cnf.Delays.ResDelay,
		Server:   cnf.Server.Port,
		Client:   cnf.Client.Port,
	}
	if config.ResDelay == 0 {
		config.ResDelay = 1
	}
	if config.ReqDelay == 0 {
		config.ReqDelay = 1
	}

	req := &tcpService.RequestForTcp{Config: config}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.StartTcpServer(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", res.GetMetric())

}
