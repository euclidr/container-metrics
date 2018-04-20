package contmetric

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readNumberFromFile(name string) (n uint64, err error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return n, err
	}

	n, err = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return n, err
	}
	return n, nil
}

func readIntFromFile(name string) (n int64, err error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return n, err
	}

	n, err = strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return n, err
	}
	return n, nil
}

func readMapFromFile(name string) (m map[string]uint64, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m = make(map[string]uint64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		v, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}
		m[parts[0]] = v
	}

	return m, nil
}
