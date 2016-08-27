package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Chane to true if you want to see the output on the console too.
var outputConsole = false

// Package-level variables
var MathUnit = new(MathDevice)   // math_device.go
var processor = new(pro)         // processor.go
var instruction = instructions() // instructions.go
var Mem []uint32
var Result []string
var newBranch = false
var initialState = true
var endState = false
var Terminal = []byte{}

// Main Function
func main() {

	// time track
	start := time.Now()

	Mem = load()

	for {
		if endState {
			break
		}

		if !newBranch {
			if !initialState {
				processor.PC++
			} else {
				processor.PC = 0
				initialState = false
			}
		} else {
			newBranch = false
		}

		processor.IR = Mem[processor.PC]
		Op := processor.IR >> 26

		switch instruction[Op].(type) {
		case TypeU:
			this := instruction[Op].(TypeU)
			this.decode(processor.IR)
			reflect.ValueOf(&this).MethodByName(this.fn).Call([]reflect.Value{})
		case TypeF:
			this := instruction[Op].(TypeF)
			this.decode(processor.IR)
			reflect.ValueOf(&this).MethodByName(this.fn).Call([]reflect.Value{})
		case TypeS:
			this := instruction[Op].(TypeS)
			this.decode(processor.IR)
			reflect.ValueOf(&this).MethodByName(this.fn).Call([]reflect.Value{})
		default:
			Result = append(Result, fmt.Sprintf("[INVALID INSTRUCTION @ 0x%.8X]\n", processor.PC<<2))
			processor.PC = 2
			newBranch = true
		}

		if MathUnit.inProgress {
			if MathUnit.clock == 3 {
				Result = append(Result, fmt.Sprintf("[HARDWARE INTERRUPTION]\n"))
				processor.makeInt()
				MathUnit.reset()
			} else {
				MathUnit.clock++
			}
		}
	}

	if outputConsole {
		for i := 0; i < len(Result); i++ {
			fmt.Printf(Result[i])
		}
		save()
	} else {
		save()
	}

	elapsed := time.Since(start)
	log.Printf("Program executed successfully and took %s.", elapsed)
	log.Println(strconv.Itoa(len(Result)) + " lines is the size of the output.")
	log.Println("see " + os.Args[2] + " to check the result out!")
}
