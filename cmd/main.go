package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/marktlinn/sysMon/internal/system"
	"github.com/marktlinn/sysMon/server"
)

// server := server.NewServer()
func main() {
	log.Println("Application starting...")

	go func() {
		for {
			cpuInfo, err := system.GetCpuStats()
			if err != nil {
				log.Println(err)
			}

			memInfo, err := system.GetMemoryStats()
			if err != nil {
				log.Println(err)
			}

			diskInfo, err := system.GetDiskStats()
			if err != nil {
				log.Println(err)
			}

			fmt.Println(cpuInfo)
			fmt.Println(memInfo)
			fmt.Println(diskInfo)

			time.Sleep(3 * time.Second)
		}
	}()
	s := server.NewServer()
	err := http.ListenAndServe(":8000", &s.Mux)
	if err != nil {
		log.Println("failed to listen to server on port ':8000'")
		os.Exit(1)
	}

}
