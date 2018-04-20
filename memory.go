package contmetric

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Doc: https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
// Reference: http://linuxperf.com/?p=142

type MemStat struct {
	Total      uint64 // min(hierarchical_memory_limit, host memory total)
	RSS        uint64 // rss in memory.stat + mapped_file
	Cached     uint64 // mapped_file + unmapped_file + tmpfs
	MappedFile uint64 // mapped_file

	SwapTotal uint64
	SwapUsed  uint64
}

func (ms *MemStat) Print() {
	var M uint64 = 1024 * 1024
	fmt.Printf("Total: %d RSS: %d Cached: %d MappedFile: %d\n", ms.Total/M, ms.RSS/M, ms.Cached/M, ms.MappedFile/M)
	fmt.Printf("SwapTotal: %d SwapUsed: %d\n", ms.SwapTotal/M, ms.SwapUsed/M)
}

func CurrentMemStat() (stat *MemStat, err error) {
	m, err := readMapFromFile("/sys/fs/cgroup/memory/memory.stat")
	stat = &MemStat{}
	stat.Total, err = totalMemory(m)
	if err != nil {
		return nil, err
	}

	stat.SwapTotal, stat.SwapUsed = swapState(m)

	stat.Cached = m["total_cache"]
	stat.MappedFile = m["total_mapped_file"]
	stat.RSS = m["total_rss"] + stat.MappedFile
	return
}

func getHostMemTotal() (n uint64, err error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != "MemTotal:" {
			continue
		}
		parts[1] = strings.TrimSpace(parts[1])
		value := strings.TrimSuffix(parts[1], "kB")
		value = strings.TrimSpace(value)
		n, err = strconv.ParseUint(value, 10, 64)
		n *= 1024
		if err != nil {
			return 0, err
		}
		break
	}
	return
}

func totalMemory(m map[string]uint64) (uint64, error) {
	hostTotal, err := getHostMemTotal()
	if err != nil {
		return 0, err
	}
	limit, ok := m["hierarchical_memory_limit"]
	if !ok {
		return 0, fmt.Errorf("missing hierarchical_memory_limit")
	}
	if hostTotal > limit {
		return limit, nil
	}
	return hostTotal, nil
}

func swapState(m map[string]uint64) (total uint64, used uint64) {
	memSwap, ok := m["hierarchical_memsw_limit"]
	if !ok {
		return 0, 0
	}

	mem := m["hierarchical_memory_limit"]
	if memSwap == mem {
		return 0, 0
	}

	total = memSwap - mem
	used = m["total_swap"]
	return total, used
}
