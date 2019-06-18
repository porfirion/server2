// Package next provides new logic implementation.
// Logic - is infinite loop, which receives some inputs and does simulation. Each simulation
// step has constant game time duration.
// There are three main type of acting:
//   1) SimulationModeContinuous - each step of simulation is driven by time (every N milliseconds).
//	 2) SimulationModeStepByStep - each step is forced by external command "simulate"
//	 3) SimulationModeReplay - same as stepbystep, but "simulate" commands are generated internally
//	    (not sure to be different from step by step, maybe should remove it)
// Inputs are passed into logic through chan. They are executed in FIFO order in single gorouting
// (for example - it is guaranteed that when you send "simulate" command -
// all previous inputs would be processed before that command).
//
package next
