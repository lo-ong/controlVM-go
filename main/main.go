package main

import (
	_ "testing"
	vbox "github.com/terra-farm/go-virtualbox"
	"fmt"

	//"time"
)

func main()  {
	//machines, err := vbox.ListMachines()
	//fmt.Println(len(machines))
	//fmt.Println(err)
	//if err != nil {
	//	os.Exit(0)
	//}
	//
	//for _, machine := range machines {
	//	fmt.Println(machine.Name)
	//}


	vbm, _ := vbox.GetMachine("bg01")
	fmt.Println(vbm.Name)

	vbm.Start()




	//time.Sleep(5 * time.Second)
	//vbm.Stop()




}