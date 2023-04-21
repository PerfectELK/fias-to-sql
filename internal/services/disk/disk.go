package disk

import (
	"github.com/ricochet2200/go-disk-usage/du"
	"os"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
)

type UsageGB struct {
	du.DiskUsage
	FreeGB uint64
	UsedGB uint64
	SizeGB uint64
}

func newUsageGB(usage du.DiskUsage) *UsageGB {
	usageGB := UsageGB{DiskUsage: usage}
	usageGB.SizeGB = usageGB.Size() * GB
	usageGB.FreeGB = usageGB.Free() * GB
	usageGB.UsedGB = usageGB.Used() * GB

	return &usageGB
}

func Usage() (*UsageGB, error) {
	pwd, err := os.Getwd()
	usageGB := &UsageGB{}
	if err != nil {
		return usageGB, err
	}
	usage := du.NewDiskUsage(pwd)
	usageGB = newUsageGB(*usage)

	return usageGB, nil
}
