package types

import "math/big"

const (
	// RemovedTailSampleRatioPpm is the percentage (in ppm) of funding samples to be removed on each
	// end of the sorted funding samples collected during a funding-tick epoch.
	// For example, if 60 funding sample entries were collected during an epoch,
	// `RemovedTailSamplePctPpm` of the top and bottom funding sample values are removed before
	// taking the average.
	// TODO(DEC-1105): Move this constant to state so that it can be changed via governance.
	RemovedTailSampleRatioPpm uint32 = 0
)

// TempBigImpactNotionalAmount is a temporary function that returns 5_000 USDC as
// the impact notional amount.
// TODO(DEC-1308): Remove this.
func TempBigImpactNotionalAmount() *big.Int {
	return big.NewInt(5_000_000_000)
}
