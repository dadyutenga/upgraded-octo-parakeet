package dashboard

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// SystemInfo holds system metrics for display.
type SystemInfo struct {
	DateTime    string
	CPUUsage    float64 // percentage (0-100)
	MemTotal    uint64  // bytes
	MemUsed     uint64  // bytes
	MemPercent  float64 // percentage (0-100)
	GoRoutines  int
	NumCPU      int
	OS          string
	Arch        string
}

// ANSI color codes for terminal output.
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

// ColorForPercent returns a color code based on the usage percentage.
func ColorForPercent(pct float64) string {
	switch {
	case pct >= 90:
		return Red
	case pct >= 70:
		return Yellow
	case pct >= 50:
		return Magenta
	default:
		return Green
	}
}

// FormatBytes converts bytes to a human-readable string.
func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// ProgressBar generates a text-based progress bar of the given width.
func ProgressBar(percent float64, width int) string {
	if width < 2 {
		width = 2
	}
	filled := int(percent / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s]", bar)
}

// GetMemInfo reads memory info from /proc/meminfo (Linux only).
func GetMemInfo() (total, used uint64, err error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		// Values in /proc/meminfo are in kB
		switch fields[0] {
		case "MemTotal:":
			memTotal = val * 1024
		case "MemAvailable:":
			memAvailable = val * 1024
		}
	}

	return memTotal, memTotal - memAvailable, nil
}

// GetCPUUsage samples /proc/stat to compute overall CPU usage over a duration.
func GetCPUUsage(sampleDuration time.Duration) (float64, error) {
	read := func() (idle, total uint64, err error) {
		file, err := os.Open("/proc/stat")
		if err != nil {
			return 0, 0, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "cpu ") {
				fields := strings.Fields(line)
				if len(fields) < 5 {
					return 0, 0, fmt.Errorf("unexpected /proc/stat format")
				}
				var vals []uint64
				for _, f := range fields[1:] {
					v, err := strconv.ParseUint(f, 10, 64)
					if err != nil {
						return 0, 0, err
					}
					vals = append(vals, v)
				}
				for _, v := range vals {
					total += v
				}
				if len(vals) >= 4 {
					idle = vals[3]
				}
				return idle, total, nil
			}
		}
		return 0, 0, fmt.Errorf("/proc/stat: cpu line not found")
	}

	idle1, total1, err := read()
	if err != nil {
		return 0, err
	}
	time.Sleep(sampleDuration)
	idle2, total2, err := read()
	if err != nil {
		return 0, err
	}

	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)
	if totalDelta == 0 {
		return 0, nil
	}

	return (1.0 - idleDelta/totalDelta) * 100.0, nil
}

// Collect gathers current system information.
func Collect(cpuSampleDuration time.Duration) SystemInfo {
	info := SystemInfo{
		DateTime:   time.Now().Format("2006-01-02 15:04:05"),
		GoRoutines: runtime.NumGoroutine(),
		NumCPU:     runtime.NumCPU(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}

	if total, used, err := GetMemInfo(); err == nil {
		info.MemTotal = total
		info.MemUsed = used
		if total > 0 {
			info.MemPercent = float64(used) / float64(total) * 100.0
		}
	}

	if cpu, err := GetCPUUsage(cpuSampleDuration); err == nil {
		info.CPUUsage = cpu
	}

	return info
}

// Render formats the SystemInfo into a colored dashboard string.
func Render(info SystemInfo) string {
	var b strings.Builder

	header := Bold + Cyan + "╔══════════════════════════════════════╗" + Reset + "\n"
	header += Bold + Cyan + "║       🖥  Dev Dashboard              ║" + Reset + "\n"
	header += Bold + Cyan + "╚══════════════════════════════════════╝" + Reset + "\n"
	b.WriteString(header)

	b.WriteString(fmt.Sprintf("\n  %s📅 Date/Time:%s  %s\n", Bold, Reset, info.DateTime))
	b.WriteString(fmt.Sprintf("  %s💻 Platform:%s   %s/%s (%d CPUs)\n", Bold, Reset, info.OS, info.Arch, info.NumCPU))

	cpuColor := ColorForPercent(info.CPUUsage)
	b.WriteString(fmt.Sprintf("\n  %s⚡ CPU Usage:%s  %s%.1f%%%s  %s\n",
		Bold, Reset, cpuColor, info.CPUUsage, Reset, ProgressBar(info.CPUUsage, 20)))

	memColor := ColorForPercent(info.MemPercent)
	b.WriteString(fmt.Sprintf("  %s🧠 Memory:%s    %s%.1f%%%s  %s  (%s / %s)\n",
		Bold, Reset, memColor, info.MemPercent, Reset,
		ProgressBar(info.MemPercent, 20),
		FormatBytes(info.MemUsed), FormatBytes(info.MemTotal)))

	return b.String()
}
