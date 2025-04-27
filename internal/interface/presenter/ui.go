package presenter

import (
	"fmt"
	"time"	
) 

func PackagingStartBanner(){
	banner := `[PACK] -------------------------------------------------------
[PACK]  🎁            STARTING PACKAGING PROCESS            🎁
[PACK] -------------------------------------------------------`
	fmt.Println(banner)
}

func DeliveyStartBanner(){
	banner := `[DELIVERY] -------------------------------------------------------
[DELIVERY]  🚚            STARTING DELIVERY PROCESS            🚚
[DELIVERY] -------------------------------------------------------`
	fmt.Println(banner)
}

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

func ShowError(err error) {
	output := fmt.Sprintf("[ERROR] ❌ %v\n", err)
	fmt.Println(output)
}

func ShowSuccess(message string){
    output := fmt.Sprintf("[SUCCESS] ✅ %s\n", message)
	fmt.Println(output)
}