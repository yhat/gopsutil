// +build !darwin,!linux,!freebsd,!windows

package disk

import "github.com/yhat/gopsutil/internal/common"

func IOCounters() (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}

func Partitions(all bool) ([]PartitionStat, error) {
	return []PartitionStat{}, common.ErrNotImplementedError
}

func Usage(path string) (*UsageStat, error) {
	return nil, common.ErrNotImplementedError
}
