package utils

import (
	"strconv"
	"strings"
)

func GetStorageCapacity(storage string) (int, error) {
	//读出storage中的数字和存储单位
	var a int = 1
	multipliers := map[string]int{
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
	}
	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(storage, suffix) {
			a = multiplier
			storage = strings.TrimSuffix(storage, suffix)
			break
		}
	}
	size, err := strconv.ParseInt(storage, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(size) * a, nil
}
