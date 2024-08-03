package main

import (
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

	s := server.NewServer()
	go func(s *server.Server) {
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

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			html := ` <td hx-swap-oob="innerHTML:#data-timestamp">` + timeStamp + `</td>
			<td hx-swap-oob="innerHTML:#cpu-data">` + cpuInfo + `</td>
			<td hx-swap-oob="innerHTML:#memory-data">` + memInfo + `</td>
			<td hx-swap-oob="innerHTML:#disk-data">` + diskInfo + `</td>
			
			`

			s.Broadcast([]byte(html))
			s.Broadcast([]byte(cpuInfo))
			s.Broadcast([]byte(memInfo))
			s.Broadcast([]byte(diskInfo))

			time.Sleep(3 * time.Second)
		}
	}(s)
	err := http.ListenAndServe(":8000", &s.Mux)
	if err != nil {
		log.Println("failed to listen to server on port ':8000'")
		os.Exit(1)
	}

}
