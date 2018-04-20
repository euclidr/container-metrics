# Container Metrics

Go package for getting docker container resource limit&usage metrics.

Tested in Ubuntu host and Alpine container.

### Usage

```
import (
    "time"
    metrics "github.com/euclidr/container-metrics"
)

metrics.GetCPUStat(time.Second, func(stat *metrics.CPUStat, err error{
    if err != nil {
        stat.Print()
    }
}))

mStat, err := metrics.CurrentMemStat()
if err != nil {
    mStat.Print()
}

dStat, err := metrics.CurrentDiskStat()
if err != nil {
    dStat.Print()
}

nStat, err := metrics.CurrentNetworkStat()
if err != nil {
    nStat.Print()
}
```