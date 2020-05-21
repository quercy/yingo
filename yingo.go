package yingo

//import "fmt"

var YIN_SAMPLING_RATE float32 = 44100

type Yin struct {
	BufferSize  int
	yinBuffer   *[]float32
	probability float32
	Threshold   float32
}

// API

func (y *Yin) YinInit(bufSize int, thresh float32) {
	y.BufferSize = bufSize
	y.Threshold = thresh

	buff := make([]float32, y.BufferSize/2)
	y.yinBuffer = &buff

}

func (y *Yin) GetPitch(d *[]float32) float32 {

	tauEstimate := -1
	var pitchInHertz float32 = -1

	y.yinDiff(d)

	y.yinCMND()

	tauEstimate = y.yinAbsThresh()

	if tauEstimate != -1 {
		pitchInHertz = YIN_SAMPLING_RATE / y.yinPI(tauEstimate)
	}

	return pitchInHertz

}

func (y *Yin) GetProb() float32 {
	return y.probability
}

//Yin private methods

// Step1: ACF

//Step2: Improving on the autocorrelation function using amplitude difference for each window at different time shifts tau

func (y *Yin) yinDiff(data *[]float32) {

	var delta float32
	if len(*data) < 10 {
		return
	}
	for tau := 0; tau < y.BufferSize/2; tau++ {
		//fmt.Println(tau)
		for i := 0; i < y.BufferSize/2; i++ {
			// fmt.Println(i, tau, len(*data))
			delta = (*data)[i] - (*data)[i+tau]

			(*y.yinBuffer)[tau] += delta * delta
		}

	}

}

//Step3: Cummulative Mean Normal Difference to deal with zero-lag errors post difference function. Set the first zero-lag difference
//       to 1 to deal with too high errors.

func (y *Yin) yinCMND() {

	var runningSum float32
	(*y.yinBuffer)[0] = 1

	for tau := 1; tau < y.BufferSize/2; tau++ {
		runningSum += (*y.yinBuffer)[tau]
		(*y.yinBuffer)[tau] *= float32(tau) / runningSum
	}
}

//Step4: Thresholding to pick the frst dip(the difference) lower than the threshold to reduce octave errors.

func (y *Yin) yinAbsThresh() int {

	var tau int
	for tau = 2; tau < y.BufferSize/2; tau++ {
		if (*y.yinBuffer)[tau] < y.Threshold {
			for tau+1 < y.BufferSize/2 && (*y.yinBuffer)[tau+1] < (*y.yinBuffer)[tau] {
				tau++
			}

			y.probability = 1 - y.Threshold
			break
		}
	}

	if tau == y.BufferSize/2 || (*y.yinBuffer)[tau] >= y.Threshold {
		tau = -1
		y.probability = 0
	}

	return tau

}

//Step5: The process is carried out for integer time-shifts (multiples of sampling rate). However, there may be a better
//       overlap at a non-integer time-shift (tau). Fit a parabolic curve to get
//       a better non-integer estimate.

func (y *Yin) yinPI(tauEstimate int) float32 {

	var betterTau float32
	var x0, x2 int

	if tauEstimate < 0 {
		x0 = tauEstimate
	} else {
		x0 = tauEstimate - 1
	}

	if tauEstimate+1 < y.BufferSize/2 {
		x2 = tauEstimate + 1
	} else {
		x2 = tauEstimate
	}

	if x0 == tauEstimate {
		if (*y.yinBuffer)[tauEstimate] <= (*y.yinBuffer)[x2] {
			betterTau = float32(tauEstimate)
		} else {
			betterTau = float32(x0)
		}
	} else if x2 == tauEstimate {
		if (*y.yinBuffer)[tauEstimate] <= (*y.yinBuffer)[x0] {
			betterTau = float32(tauEstimate)
		} else {
			betterTau = float32(x0)
		}
	} else {
		var s0, s1, s2 float32
		s0 = (*y.yinBuffer)[x0]
		s1 = (*y.yinBuffer)[tauEstimate]
		s2 = (*y.yinBuffer)[x2]

		betterTau = float32(tauEstimate) + (s2-s0)/(2*(2*s1-s2-s0))

	}

	return betterTau
}
