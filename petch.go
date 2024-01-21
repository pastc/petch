package main

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

const (
	space = 8
)

func spaces(key string) string {
	return key + strings.Repeat(" ", space-len(key))
}

func getTitle(h types.Host) string {
	user, err := user.Current()
	if err != nil {
		user.Username = ""
	}
	return user.Username + "@" + h.Info().OS.Platform
}

func getMem(h types.Host) string {
	mem, err := h.Memory()
	if err != nil {
		mem.Used, mem.Total = 0, 0
	}
	return formatBytes(mem.Used) + "/" + formatBytes(mem.Total)
}

func formatTime(duration time.Duration) string {
	time := int(duration.Abs().Seconds())
	var (
		days    int
		hours   int
		minutes int
		// seconds int
	)
	days = time / 86400
	hours = (time - (days * 86400)) / 3600
	minutes = (time - (days * 86400) - (hours * 3600)) / 60
	// seconds = time - (days * 86400) - (hours * 3600) - (minutes * 60)

	return strconv.Itoa(days) + "d " + strconv.Itoa(hours) + "h " + strconv.Itoa(minutes) + "m"
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	host, err := sysinfo.Host()
	if err != nil {
		panic(err)
	}
	info := host.Info()
	fmt.Println()
	fmt.Println(getTitle(host))
	fmt.Println(strings.Repeat("-", len(getTitle(host))))
	fmt.Println(spaces("OS"), info.OS.Name, info.OS.Version, info.Architecture)
	fmt.Println(spaces("Kernel"), info.KernelVersion)
	fmt.Println(spaces("Timezone"), info.Timezone, info.TimezoneOffsetSec)
	fmt.Println(spaces("Uptime"), formatTime(info.Uptime()))
	fmt.Println(spaces("Memory"), getMem(host))
	fmt.Println()
}
