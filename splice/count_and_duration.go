package splice

import "time"

func getCountAndDuration(timeoutSec uint) (uint, time.Duration) {
	var timeoutCnt uint
	var timeoutDuration time.Duration

	if timeoutSec == 0 {
		return 1, time.Hour * 24 * 30
	}

	if timeoutSec > 0 {
		if timeoutSec <= 20 {
			timeoutCnt = timeoutSec
			timeoutDuration = time.Second
		} else if timeoutSec <= 100 {
			timeoutCnt = timeoutSec / 5
			if timeoutSec%5 > 0 {
				timeoutCnt++
			}
			timeoutDuration = 5 * time.Second
		} else if timeoutSec <= 200 {
			timeoutCnt = timeoutSec / 10
			if timeoutSec%10 > 0 {
				timeoutCnt++
			}
			timeoutDuration = 10 * time.Second
		} else {
			timeoutCnt = timeoutSec / 15
			if timeoutSec%15 > 0 {
				timeoutCnt++
			}
			timeoutDuration = 15 * time.Second
		}
	}

	return timeoutCnt, timeoutDuration
}
