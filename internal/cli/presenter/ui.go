package presenter

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"time"
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

func ShowError(err error) {
	output := fmt.Sprintf("%s[ERROR]%s: %v", ColorRed, ColorReset, err)
	fmt.Println(output)
}

func ShowSuccess(message string) {
	output := fmt.Sprintf("%s[SUCCESS]%s %s", ColorGreen, ColorReset, message)
	fmt.Println(output)
}

func ShowInfo(message string) {
	output := fmt.Sprintf("%s[INFO]%s %s", ColorBlue, ColorReset, message)
	fmt.Println(output)
}

func ShowMessage(message string) {
	output := fmt.Sprintf("%s", message)
	fmt.Println(output)
}

func ShowMessageSpinner(message string) *spinner.Spinner {
	output := fmt.Sprintf("%s 🔨 : ", message)
	sp := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
	sp.Prefix = output
	sp.Start()
	return sp
}

func Show(response result.DomainResult, fileLogger *logger.FileLogger) {
	if response.IsSuccess() {
		ShowSuccess(response.Message)
	} else {
		ShowError(response.Error)
		showFileLogger(fileLogger)
	}
}

func showFileLogger(fileLogger *logger.FileLogger) {
	logs, err := fileLogger.ReadFromFile()
	if err != nil {
		ShowError(err)
	}
	for _, log := range logs {
		if log.Level == logger.ERROR {
			ShowError(log.Error)
		}
		if log.Level == logger.SUCCESS {
			ShowSuccess(log.Message)
		}
		if log.Level == logger.INFO {
			ShowInfo(log.Message)
		} else {
			ShowMessage(log.Message)
		}
	}
}
