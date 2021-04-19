package modt


import (
	"log"
	"os/exec"
)

func SendGcodeTemp(filename string) {
	log.Println("Sending gcode.....\n\n\n\n")
	cmd := exec.Command("python", "send_gcode.py", filename)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
	
}