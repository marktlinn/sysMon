package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marktlinn/sysMon/internal/system"
)

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
	time.Sleep(5 * time.Minute)
}
