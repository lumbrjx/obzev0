package main

import (
	"fmt"
	"obzev0/common/definitions"
	"os"

	"gopkg.in/yaml.v2"
)

func generateYaml(serverAddr string, path string) {
	config := definitions.Config{
		ServerAddr: serverAddr,
		LatencySvcConfig: definitions.LatencySvcConfig{
			Enabled:  true,
			ReqDelay: 1,
			ResDelay: 1,
			Server:   "7070",
			Client:   "127.0.0.1:8080",
		},
		TcAnalyserSvcConfig: definitions.TcAnalyserSvcConfig{
			Enabled:  false,
			NetIFace: "eth0",
		},
		PacketManipulationSvcConfig: definitions.PacketManipulationSvcConfig{
			Enabled:         false,
			Server:          "9091",
			Client:          "127.0.0.1:8080",
			DropRate:        "0.8",
			CorruptRate:     "0.4",
			DurationSeconds: 8,
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error marshalling struct to YAML:", err)
		return
	}

	file, err := os.Create(path + "obzev0cnf.yaml")
	if err != nil {
		fmt.Println("Error creating YAML file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to YAML file:", err)
		return
	}

	fmt.Println("YAML file created successfully.")
}
