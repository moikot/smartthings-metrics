package extractors

var (
	motionSensorValues = map[string]float64{
		"inactive": 0.0,
		"active":   1.0,
	}

	switchValues = map[string]float64{
		"off": 0.0,
		"on":  1.0,
	}

	waterSensorValues = map[string]float64{
		"dry": 0.0,
		"wet": 1.0,
	}
)
