////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                           //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package exponential

// Default constants returned by DefaultMovingAvgParams.
const (
	defaultCutoff          = 0.5
	defaultInitialAverage  = 0.15
	defaultSmoothingFactor = 2
	defaultNumberOfEvents  = 100
)

// MovingAvgParams are the parameters needed to create a new MovingAvg.
type MovingAvgParams struct {
	// Cutoff is the maximum the moving average can reach before an error is
	// returned on intake. The value should range from 0.0 to 1.0.
	Cutoff float32

	// InitialAverage is the initial exponential moving average to start with.
	InitialAverage float32

	// SmoothingFactor, the exponential smoothing factor; gives the most recent
	// events more weight. The greater the smoothing factor, the greater the
	// influence of more recent events.
	SmoothingFactor float32

	// NumberOfEvents is the number of events to average over.
	NumberOfEvents uint32
}

// DefaultMovingAvgParams returns MovingAvgParams with default values.
func DefaultMovingAvgParams() MovingAvgParams {
	return MovingAvgParams{
		Cutoff:          defaultCutoff,
		InitialAverage:  defaultInitialAverage,
		SmoothingFactor: defaultSmoothingFactor,
		NumberOfEvents:  defaultNumberOfEvents,
	}
}
