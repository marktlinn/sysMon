package system

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetCpuStats() (string, error) {
	cpuStats, err := cpu.Info()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve CPU stats: %w", err)
	}

	cpuOutput := fmt.Sprintf("CPU: %s\nCores: %d", cpuStats[0].ModelName, len(cpuStats))
	return cpuOutput, nil
}

func GetMemoryStats() (string, error) {
	sysOS := runtime.GOOS

	vmStats, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve VirtualMemory stats: %w", err)
	}

	hostStats, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Host stats: %w", err)
	}

	memUsedGB := convertBytesToGB(vmStats.Used)
	memTtlGB := convertBytesToGB(vmStats.Total)

	memOutput := fmt.Sprintf(
		"Hostname: %s\nTotal Memory %.2f GB\nUsed Memory: %.2f GB\nOS: %s",
		hostStats.Hostname,
		memTtlGB,
		memUsedGB,
		sysOS,
	)

	return memOutput, nil
}

func GetDiskStats() (string, error) {
	diskStats, err := disk.Usage("/")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Disk stats: %w", err)
	}

	diskTtlGB := convertBytesToGB(diskStats.Total)
	diskFreeGB := convertBytesToGB(diskStats.Free)

	diskOutput := fmt.Sprintf(
		"Total Disk Space: %.2f GB\nFree Disk Space: %.2f GB",
		diskTtlGB,
		diskFreeGB,
	)

	return diskOutput, nil
}

func convertBytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
