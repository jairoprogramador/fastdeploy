package presenter

import (
	"deploy/internal/domain/model"
	"fmt"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func ShowBanner() {
	banner := `
    ____  _____ ____  _     _____   __
   |  _ \| ____|  _ \| |   / _ \ \ / /
   | | | |  _| | |_) | |  | | | \ V / 
   | |_| | |___|  __/| |__| |_| || |  
   |____/|_____|_|   |_____\___/ |_|  
   :: CLI Aplication ::        (v1.0.0)
   `
	fmt.Println(banner)
}

func ShowStart(message string) {
	output := fmt.Sprintf("%s[START] 🚚%s %s", ColorPurple, ColorReset, message)
	fmt.Println(output)
}

func ShowError(stepName string, err error) {
	output := fmt.Sprintf("%s[ERROR]%s %s: %v", ColorRed, ColorReset, stepName, err)
	fmt.Println(output)
}

func ShowSuccess(stepName string, message string) {
	output := fmt.Sprintf("%s[SUCCESS]%s %s: %s", ColorGreen, ColorReset, stepName, message)
	fmt.Println(output)
}

func ShowLogStore(logStore *model.LogStore) {
	for _, step := range logStore.Steps {
		if step.Status == model.StatusError {
			ShowError(step.StepName, step.Error)
			return
		}
	}

	for _, step := range logStore.Steps {
		ShowSuccess(step.StepName, step.Message)
	}
}
