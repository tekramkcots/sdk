package indicators

func SimpleMovingAverage(noOfDays int, values []float64) []float64 {
	var sma []float64
	for i := 0; i < len(values); i++ {
		if i < noOfDays {
			sma = append(sma, 0)
			continue
		}
		var sum float64
		for j := i - noOfDays; j < i; j++ {
			sum += values[j]
		}
		sma = append(sma, sum/float64(noOfDays))
	}
	return sma
}

func ExponentialMovingAverage(noOfDays int, values []float64) []float64 {
	var ema []float64
	for i := 0; i < len(values); i++ {
		if i < noOfDays {
			ema = append(ema, 0)
			continue
		}
		if i == noOfDays {
			var sum float64
			for j := i - noOfDays; j < i; j++ {
				sum += values[j]
			}
			ema = append(ema, sum/float64(noOfDays))
			continue
		}
		ema = append(ema, (values[i]*2/(float64(noOfDays)+1))+(ema[i-1]*(float64(noOfDays)-1)/(float64(noOfDays)+1)))
	}
	return ema
}

func MovingAverageConvergenceDivergence(fastManPeriod, slowMaPeriod int, values []float64) []float64 {
	var macd []float64
	ema1 := ExponentialMovingAverage(fastManPeriod, values)
	ema2 := ExponentialMovingAverage(slowMaPeriod, values)
	for i := 0; i < len(values); i++ {
		if i < slowMaPeriod {
			macd = append(macd, 0)
			continue
		}
		macd = append(macd, ema1[i]-ema2[i])
	}
	return macd
}

func MovingAverageConvergenceDivergenceSignalLine(signalPeriod int, macd []float64) []float64 {
	var signal []float64
	for i := 0; i < len(macd); i++ {
		if i < signalPeriod {
			signal = append(signal, 0)
			continue
		}
		var sum float64
		for j := i - signalPeriod; j < i; j++ {
			sum += macd[j]
		}
		signal = append(signal, sum/float64(signalPeriod))
	}
	return signal
}
