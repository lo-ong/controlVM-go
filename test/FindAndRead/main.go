package main

import (
	"github.com/lo_ong/controlVM-go/engine"
	"fmt"
)

const filePath string = "../../source/operation.json"

func main() {
	vms, rvmNames, _ := engine.FindMachineInLocal()
	//for _, vm := range vms{
	//	fmt.Println(vm.Name)
	//}
	data, _ := engine.ReadJson(filePath)



	virtualMachines, err := engine.BuildMachine(vms, rvmNames, data)
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
