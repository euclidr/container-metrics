package contmetric

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DiskStat struct {
	Read  uint64
	Write uint64
}

func (ds *DiskStat) Print() {
	fmt.Printf("Read: %d, Write: %d\n", ds.Read, ds.Write)
}

func CurrentDiskStat() (stat *DiskStat, err error) {
	var read, write uint64
	if diskAcctFile == "" {
		for _, file := range diskAcctFiles {
			read, write, _ = getDiskReadWrite(file)
			if read+write > 0 {
				diskAcctFile = file
				break
			}
		}
	} else {
		read, write, err = getDiskReadWrite(diskAcctFile)
	}

	if err != nil {
		return nil, err
	}

	return &DiskStat{Read: read, Write: write}, nil
}

var diskAcctFile string

var diskAcctFiles = []string{
	"/sys/fs/cgroup/blkio/blkio.io_service_bytes_recursive",
	"/sys/fs/cgroup/blkio/blkio.throttle.io_service_bytes"}

func getDiskReadWrite(name string) (read, write uint64, err error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var r, w uint64
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}
		if parts[1] == "Read" {
			tmp, _ := strconv.Atoi(parts[2])
			r += uint64(tmp)
			continue
		}
		if parts[1] == "Write" {
			tmp, _ := strconv.Atoi(parts[2])
			w += uint64(tmp)
			continue
		}
	}
	return r, w, nil
}
