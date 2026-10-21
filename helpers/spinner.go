package helpers

import (
	"fmt"
	"time"
)

// showSpinner prints a simple animated spinner while scanning is in progress.
func (f *FileStats) ShowSpinner(stop <-chan struct{}) {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	dots := []string{"...", ".", ".."}

	ticker := time.NewTicker(time.Microsecond * 80)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-stop:
			fmt.Print("\r\033[K") // clear the spinner line
			return
		case <-ticker.C:
			fmt.Printf("\r%s Scanning project%s ", spinners[i%len(spinners)], dots[i%len(dots)])
			i++
		}
	}
}
