package main

import (
	"github.com/lo_ong/controlVM-go/engine"
	"fmt"
	"time"
	"strings"
)

const filePath string = "../../source/operation.json"

func main() {
	// print cost time
	start := time.Now()
	defer func() {
		cost := time.Since(start)
		fmt.Println("App finish cost: ", cost)
	}()

	fmt.Println(strings.Split(`"bg01"`, `"`))


	listVMNames, runningVMNames, _ := engine.FindMachineInLocal()
	fmt.Println("listVMNames:", listVMNames)
	fmt.Println("runningVMNames:", runningVMNames)


	data, _ := engine.ReadJson(filePath)

	virtualMachines, err := engine.BuildMachine(listVMNames, runningVMNames, data)
	if err != nil {
		fmt.Println(err)
	}

	for _, machine := range virtualMachines {
		machine.BuildCommand()
	}

	for _, machine := range virtualMachines {
		machine.Run()
	}
}
