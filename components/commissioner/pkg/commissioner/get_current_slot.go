package commissioner

import (
	"fmt"
)

func getCurrentSlot() (int, error) {
	tip, err := queryTip()
	if err != nil {
		return 0, fmt.Errorf("queryTip: %v", err)
	}
	return tip.Slot, nil
}
