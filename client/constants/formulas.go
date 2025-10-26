package constants

import "math"

const (
	NPC_DAMAGE_VARIATION_PCT = 0.12 // 12%
)

func XP_CAP(level uint64) uint64 {
	return uint64(float64(level)*100*1.25 + math.Pow(float64(level), 2.0)*66.0)
}
