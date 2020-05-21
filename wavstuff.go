package yingo

import (
	"log"
	"math"
	"os"

	"github.com/go-audio/wav"
)

// Pitch ...
type Pitch struct {
	HopStamp         int
	Detectedpitch    float32
	PitchProbability float32
	StdFrequency     float32
	MidiNumber       int
}

// MonoAnalyser ...
func MonoAnalyser(f string, bufferapproximate bool, hopSize int) <-chan Pitch {
	// Open the file
	rawFile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer rawFile.Close()

	// Create a decoder for the file
	wavDecoder := wav.NewDecoder(rawFile)
	pcm, err := wavDecoder.FullPCMBuffer()
	if err != nil {
		log.Fatal(err)
	}

	// Maps note names -> hz values
	freqArray := loadFreqArray()

	// Load the data from the wave file into memory directly
	intBuffer := pcm.AsIntBuffer().Data
	pcmArray := make([]float32, len(intBuffer))

	// Yin alg expects float32s
	for i, val := range intBuffer {
		pcmArray[i] = float32(val)
	}

	// Get the number of iterations we need to do
	iterations := len(pcmArray) / hopSize
	pch := make(chan Pitch, iterations)

	for i := 0; i < iterations; i++ {
		// Init yin, create variables
		yin := Yin{}
		yin.YinInit(hopSize, float32(0.05))
		batch := make([]float32, hopSize)
		batch = pcmArray[i*hopSize : (i*hopSize + hopSize)]
		pitch := Pitch{HopStamp: i}

		// Do the Yin alg computations, add to our data structure
		pitch.Detectedpitch = yin.GetPitch(&batch)
		pitch.PitchProbability = yin.GetProb()
		pitch.StdFrequency, pitch.MidiNumber = moarData(pitch.Detectedpitch, freqArray)

		pch <- pitch // Send to the channel
	}

	close(pch)
	return pch
}

func loadFreqArray() *[88]float32 {
	// Loop over all MIDI values
	var frArr [88]float32
	for i := 21; i <= 108; i++ {
		//midi to frequency; A4 = 440 Hz

		//very dirty; could use 1 << (exponent)or 1 >> (exponent) after casting to unsigned ints; not decided
		frArr[i-21] = float32(math.Pow(2, float64((i-69))/float64(12.0)) * 440)

	}
	return &frArr
}

func basicAbs(x float32) float32 {

	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}

func moarData(p float32, freqArray *[88]float32) (float32, int) {

	pitch := p

	if pitch == -1 {
		return 0, 0
	}

	smallestDiff := basicAbs(pitch - (*freqArray)[0])
	var stdFrequency float32
	var midiNumber int
	for n, val := range *freqArray {
		xDiff := basicAbs(pitch - val)
		if xDiff < smallestDiff {
			smallestDiff = xDiff
			stdFrequency = val
			midiNumber = n + 21
		}
	}

	return stdFrequency, midiNumber

}
