package contmetric

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var ethInterface string

var (
	ErrDefaultEthInterfaceNotfound = fmt.Errorf("default EthInterface notfound")
)

type NetworkStat struct {
	RxBytes uint64
	TxBytes uint64
}

func (ns *NetworkStat) Print() {
	fmt.Printf("RxBytes: %d, TxBytes: %d\n", ns.RxBytes, ns.TxBytes)
}

func CurrentNetworkStat() (stat *NetworkStat, err error) {
	if ethInterface == "" {
		return nil, ErrDefaultEthInterfaceNotfound
	}
	folder := "/sys/class/net/" + ethInterface + "/statistics/"
	r, err := readNumberFromFile(folder + "rx_bytes")
	if err != nil {
		return nil, err
	}
	t, err := readNumberFromFile(folder + "tx_bytes")
	if err != nil {
		return nil, err
	}

	return &NetworkStat{RxBytes: r, TxBytes: t}, nil
}

func init() {
	// $ ip -o -4 route show to default
	// default via 172.17.0.1 dev eth0
	cmd := exec.Command("ip", "-o", "-4", "route", "show", "to", "default")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ip cmd err: " + err.Error())
		fmt.Println("ip cmd err result: " + stderr.String())
		return
	}
	parts := strings.Split(strings.TrimSpace(out.String()), " ")
	if len(parts) != 5 {
		fmt.Println(fmt.Errorf("invalid result from \"ip -o -4 route show to default\": %s", out.String()))
		return
	}
	ethInterface = strings.TrimSpace(parts[4])
}
