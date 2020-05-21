/*
 * Uses the microphone to record chunks and runs the YIN Pitch Detection on the chunks.
 * Writes the output to a channel of floats
 * Ctrl+C to break the audio recording loop
 *
 * Requires portaudio for go : go get github.com/gordonklaus/portaudio

 */

package main

import (
	"fmt"

	"github.com/quercy/yingo"
)

func main() {

	pitchChannel := make(chan float32)
	yingo.MicInput(2048, true, &pitchChannel)

	for {
		select {
		case pitch := <-pitchChannel:
			fmt.Println("Pitch is: ", pitch)
		}
	}
}
