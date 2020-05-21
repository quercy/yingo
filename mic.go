package yingo

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

//type Mic struct {
//chunkSize int
//bufferSize int
//bufferoptimize bool
//}

func MicInput(chunkSize int, bfropt bool, pch *chan float32) {
	//yin variables
	//m = Mic{chunkSize: chnkSze, bufferSize: 100, bufferoptimize: bfropt }
	threshold := float32(0.05)
	bufferSize := 100

	//var pitch float32

	var bufferincrement int

	//pch := make(chan float32)

	if bfropt {
		bufferincrement = 1
	} else {
		bufferincrement = 100
	}

	fmt.Println("Enter loop")
	go func() {

		portaudio.Initialize()
		defer portaudio.Terminate()

		input := make([]int16, chunkSize)
		stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(input), input)
		chkErr(err)
		defer stream.Close()

		chkErr(stream.Start())
		for {
			yin := Yin{}
			fmt.Println("New Chunk")
			chkErr(stream.Read())
			fmt.Println(input)
			figurePitch(&yin, bufferincrement, input, pch, bufferSize, threshold)
		}
	}()
}

func figurePitch(yin *Yin, bufferincrement int, input []int16, pch *chan float32, bufferSize int, threshold float32) {
	var pitch float32
	fmt.Println("Processing")
	for pitch < 10 {
		//fmt.Println(bufferSize)
		if bufferSize >= len(input) {
			fmt.Println("Break")
			pitch = -1
			break
		}
		yin.YinInit(bufferSize, threshold)
		// pitch = yin.GetPitch(&input)
		bufferSize += bufferincrement
	}
	//fmt.Println(bufferSize)
	fmt.Println(pitch)
	*pch <- pitch
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}
