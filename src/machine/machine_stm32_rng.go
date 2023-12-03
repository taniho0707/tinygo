//go:build stm32 && !(stm32f103 || stm32l0x1 || stm32f446)

package machine

import "device/stm32"

var rngInitDone = false

const RNG_MAX_READ_RETRIES = 1000

// GetRNG returns 32 bits of cryptographically secure random data
func GetRNG() (uint32, error) {
	if !rngInitDone {
		initRNG()
		rngInitDone = true
	}

	if stm32.RNG.SR.HasBits(stm32.RNG_SR_CECS) {
		return 0, ErrClockRNG
	}
	if stm32.RNG.SR.HasBits(stm32.RNG_SR_SECS) {
		return 0, ErrSeedRNG
	}

	cnt := RNG_MAX_READ_RETRIES
	for !stm32.RNG.SR.HasBits(stm32.RNG_SR_DRDY) {
		cnt--
		if cnt == 0 {
			return 0, ErrTimeoutRNG
		}
	}

	ret := stm32.RNG.DR.Get()
	return ret, nil
}
