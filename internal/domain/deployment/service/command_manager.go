package service

import (
	"fmt"
	"slices"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/chain/command"
)

var stepsOrder = []string {
	constant.StepTest,
	constant.StepSupply,
	constant.StepPackage,
	constant.StepDeploy,
}

type CommandManager interface {
	CreateChain(finalStep string, blockedSteps []string) (chain.CommandChain, error)
}

type CommandManagerImpl struct {
	stepCommands map[string]chain.CommandChain
}

func NewCommandManager(factoryStrategy port.FactoryStrategy) CommandManager {
	steps := make(map[string]chain.CommandChain)
	steps[constant.StepTest] = command.NewTestCommand(factoryStrategy.CreateTestStrategy())
	steps[constant.StepSupply] = command.NewSupplyCommand(factoryStrategy.CreateSupplyStrategy())
	steps[constant.StepPackage] = command.NewPackageCommand(factoryStrategy.CreatePackageStrategy())
	steps[constant.StepDeploy] = command.NewDeployCommand(factoryStrategy.CreateDeployStrategy())

	return &CommandManagerImpl{
		stepCommands: steps,
	}
}

func (f *CommandManagerImpl) CreateChain(finalStep string, blockedSteps []string) (chain.CommandChain, error){
	var firstCommand chain.CommandChain
	var lastCommand chain.CommandChain

	addCommand := func(c chain.CommandChain) {
		if firstCommand == nil {
			firstCommand = c
			lastCommand = c
		} else {
			lastCommand.SetNext(c)
			lastCommand = c
		}
	}

	for _, stepName := range stepsOrder {
		if !slices.Contains(blockedSteps, stepName)  {
			cmd, ok := f.stepCommands[stepName]
			if !ok {
				return nil, fmt.Errorf("comando no encontrado en el mapa: %s", stepName)
			}
			addCommand(cmd)
		}
	}

	return firstCommand, nil
}
