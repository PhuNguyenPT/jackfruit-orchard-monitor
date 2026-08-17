package model

// SoilCalibration holds a single sensor's dry/wet raw ADC reference points,
// used to convert a raw reading into a moisture percentage.
type SoilCalibration struct {
	Dry int
	Wet int
}
