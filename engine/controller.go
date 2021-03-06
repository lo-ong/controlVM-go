package engine

import (
	"io/ioutil"
	"errors"
	"encoding/json"
	"strings"
	"fmt"
	"regexp"
	"os/exec"
	"time"
)

const (
	ExecRetryCount int = 100
	ExecRetrySleep time.Duration = 1*time.Second
)

// find all virtual machine in local virtual box
func FindMachineInLocal() ([]string, []string,  error) {
	// Find all virtual machine
	var listVMNames []string
	var runningVMNames []string
	cmdString := "VBoxManage list vms"
	listVMString, err := Execution(cmdString)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(listVMString, "\n")
	for _, line := range lines {
		if len(line) != 0 {
			listVMNames = append(listVMNames, strings.Split(strings.Split(line, " ")[0], `"`)[1])
		}
	}

	// Find machine is running
	cmdString = "VBoxManage list runningvms"
	runningVMString, err := Execution(cmdString)
	if err != nil {
		return nil, nil, err
	}
	lines = strings.Split(runningVMString, "\n")
	for _, line := range lines {
		if len(line) != 0 {
			runningVMNames = append(runningVMNames, strings.Split(strings.Split(line, " ")[0], `"`)[1])
		}
	}

	return listVMNames, runningVMNames, nil
}

// Read operation file source
func ReadJson(filePath string) (map[string]interface{}, error) {
	//filePath := "github.com/control-vbox/source/operation.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("read operation file failed: " + err.Error())
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	return m, nil
}

type VirtualMachineConfigure struct {
	UserName string
	Password string
	StartType string
	PowerOff bool
	Shell []string
}

type VirtualMachine struct {
	Name string
	HaveStart bool
	Configure VirtualMachineConfigure
	Commands []string
}

type VirtualMachines struct {
	vms []VirtualMachine
}

func NewVirtualMachine(name string, haveStart bool, configure VirtualMachineConfigure) *VirtualMachine {
	return &VirtualMachine{
		Name: name,
		HaveStart: haveStart,
		Configure: configure,
	}
	return nil
}

// Parse json and build machine
func BuildMachine(listVMNames []string, runningVMNames []string, inputJson map[string]interface{}) ([]*VirtualMachine, error) {
	var virtualMachines []*VirtualMachine

	for name, configure := range inputJson {
		machineExist := false
		for _, vmName := range listVMNames {
			if strings.Compare(vmName, name) != 0 {
				continue
			} else {
				machineExist = true
				break
			}
		}

		if !machineExist {
			return nil, errors.New("virtual machine " + name + " hasn't be created")
		}

		// The virtual machine is here
		// Whether the machine is running
		haveStart := false
		for _, runningName := range runningVMNames{
			if strings.Compare(runningName, name) != 0 {
				continue
			} else {
				haveStart = true
				break
			}
		}

		bytes, _ := json.Marshal(configure)
		vmc := VirtualMachineConfigure{}
		err := json.Unmarshal(bytes, &vmc)
		if err != nil {
			return nil, err
		}

		virtualMachine := NewVirtualMachine(name, haveStart, vmc)
		virtualMachines = append(virtualMachines, virtualMachine)

	}

	return virtualMachines, nil
}

func DeleteExtraSpace(s string) string {
	strNoTab := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(strNoTab, " ")
}

//
func (vm *VirtualMachine)BuildCommand() error {
	virtualName := vm.Name
	startType := vm.Configure.StartType
	userName := vm.Configure.UserName
	password := vm.Configure.Password

	if !vm.HaveStart {
		cmdString := "VBoxManage startvm "+ virtualName +" --type " + startType
		vm.AddCommand(cmdString)
	}

	shell := vm.Configure.Shell
	for _, cmd := range shell {
		cmdNoExtraSpace := DeleteExtraSpace(strings.Trim(cmd, ""))

		cmdPart := strings.Split(cmdNoExtraSpace, " ")
		if len(cmdPart) <= 0 {
			return errors.New("an empty command")
		}

		operation := cmdPart[0]
		cmdPart = append(cmdPart[1:])
		arguments := strings.Join(cmdPart, " ")

		cmdString := "VBoxManage --nologo guestcontrol "+ virtualName + " " +
			"run --exe /bin/"+ operation +" --username "+ userName +" " +
			"--password "+ password +" --wait-stdout --wait-stderr -- -l " + arguments

		vm.AddCommand(cmdString)
	}

	if vm.Configure.PowerOff {
		cmdString := "VBoxManage controlvm "+ virtualName +" acpipowerbutton"
		vm.AddCommand(cmdString)
	}

	return nil
}

func (vm *VirtualMachine)AddCommand(command string) {
	vm.Commands = append(vm.Commands, command)
}

// Open virtual machine, and run shell command
func (vm *VirtualMachine)Run() error {
	fmt.Println("run command start: ")
	defer func() {
		fmt.Println("run command finish: ")
	}()
	for _, command := range vm.Commands {
		fmt.Println("run command try: ", command)
		err := Retry(ExecRetryCount, ExecRetrySleep, func() error {
			return Executor(command)
		})

		if err != nil {
			fmt.Println("run command error :", command)
			return err
		}
		fmt.Println("run command success :", command)

	}

	return nil
}

func Executor(command string) error {
	cmd := exec.Command("bash", "-c", command)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func Execution(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			fmt.Println("error: ", err, "\n try again: ", attempts)
			time.Sleep(sleep)
			return Retry(attempts, sleep, fn)
		}
		fmt.Println("error: ", err, "\n retry fail: cancel control ")
		return err
	}
	return nil
}