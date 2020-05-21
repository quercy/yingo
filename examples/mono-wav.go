/* Simple analysis on a monophonic single frequency WAV file
 * Returns a channel of type Pitch.
 * The type Pitch carries the Hop number at which the analysis is made, the pitch detected at that Hop, the confidence of the detection.
 * Based on the detected pitch, the closest standard tuning(A4:440 Hz) frequency and the standard midi number are alo returned.
 * pitch of -1 is returned in case of silence or failure to detect the time period by the algorithm.
 * Pitch is analysed for every chunk of size defined by the third parameter, the hopSize
 * The hopSize is time-frequency trade-off. For better time resolution, use a smaller chunk size.
 * However, this would limit the detection of lower frequencies.
 *
 * Theoretically, you would want 2*SamplingRate/FrequencyToBeDetected number of samples for analysis.
 * So, with 2048 chunk size, you are looking at about 43 Hz lowest possible frequency.
 */

package main

import (
	"fmt"
	"yingo"
)

func main() {

	fmt.Println("bach.wav:")
	bachPitches := yingo.MonoAnalyser("bach.wav", true, 2048)
	printPitches(bachPitches)

	fmt.Println("\npiano.wav:")
	pianoPitches := yingo.MonoAnalyser("piano.wav", true, 2048)
	printPitches(pianoPitches)
}

func printPitches(pitches <-chan yingo.Pitch) {
	for pitch := range pitches {
		if pitch.Detectedpitch > 0 {
			fmt.Printf("hz value: %v estimated note: %v\n", pitch.Detectedpitch, midiNoteToNoteName(pitch.MidiNumber))
		}
	}
}

func midiNoteToNoteName(midiNumber int) string {
	noteNames := []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
	if midiNumber < 21 {
		return ""
	}
	return noteNames[(midiNumber-21)%len(noteNames)]
}
