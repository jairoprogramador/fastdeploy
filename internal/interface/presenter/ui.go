package presenter

import (
	"fmt"
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

func ShowLoader(done chan bool) {
    loader := []string{"|", "/", "-", "\\"}
    i := 0
    for {
        select {
        case <-done:
            return
        default:
            fmt.Printf("\rRunning... %s", loader[i%len(loader)])
            i++
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func ShowBanner(){
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

func ShowMessage(message string) {
    fmt.Println(message)
}

func ShowStart(message string) {
	output := fmt.Sprintf("%s[START] 🚚%s %s", ColorPurple, ColorReset, message)
	fmt.Println(output)
}

func ShowError(err error) {
	output := fmt.Sprintf("%s ❌  [ERROR]%s %v", ColorRed, ColorReset, err)
	fmt.Println(output)
}

func ShowSuccess(message string){
    output := fmt.Sprintf("%s ✅  [SUCCESS]%s %s", ColorGreen, ColorReset, message)
	fmt.Println(output)
}