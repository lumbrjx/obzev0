package main

import (
	"context"
	"fmt"
	"log"
	"obzev0/common/definitions"
	"obzev0/common/proto/tcAnalyser"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

func LoadConfig(filename string) (definitions.Config, error) {
	var config definitions.Config

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
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// c := latency.NewLatencyServiceClient(conn)
	t := tcAnalyser.NewTcAnalyserServiceClient(conn)

	// cnf, err := LoadConfig("obzevConf.yaml")
	// config := &latency.TcpConfig{
	// 	ReqDelay: cnf.Delays.ReqDelay,
	// 	ResDelay: cnf.Delays.ResDelay,
	// 	Server:   cnf.Server.Port,
	// 	Client:   cnf.Client.Port,
	// }
	// println(
	// 	config.Client,
	// 	config.Server,
	// 	config.ResDelay,
	// 	config.ReqDelay,
	// )

	req2 := &tcAnalyser.RequestForUserSpace{
		Config: &tcAnalyser.TcConfig{Interface: "eth0"},
	}
	// req := &latency.RequestForTcp{Config: &la}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := t.StartUserSpace(ctx, req2)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", res.GetMessage())

}
