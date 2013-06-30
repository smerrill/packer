package vmware

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"log"
	"fmt"
)

type stepProvision struct{}

func (*stepProvision) Run(state map[string]interface{}) multistep.StepAction {
	comm := state["communicator"].(packer.Communicator)
	hook := state["hook"].(packer.Hook)
	ui := state["ui"].(packer.Ui)
	driver := state["driver"].(Driver)
	vmxPath := state["vmx_path"].(string)

	ui.Say("Mounting the VMware tools")
	if err := driver.MountTools(vmxPath); err != nil {
		err := fmt.Errorf("Error mounting the VMWare Tools on VM: %s", err)
		state["error"] = err
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Println("Running the provision hook")
	if err := hook.Run(packer.HookProvision, ui, comm, nil); err != nil {
		state["error"] = err
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (*stepProvision) Cleanup(map[string]interface{}) {}
