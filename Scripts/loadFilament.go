package modt

import (
	"log"
	"os/exec"
)

//Runs the python version of this until I can figure out how configurations work
func LoadFilamentTemp() {
	log.Println("Begginging load filament\n\n\n\n\n\n")
	cmd := exec.Command("python", "load_filament.py")
	log.Printf("Running command and waiting for it to finish...")
	out,err := cmd.CombinedOutput()
	log.Printf("Command finished with error: %v", err)
	log.Println("Output:\n"+string(out))
	log.Println("ending load filament\n\n\n\n\n\n")
}

const (
	loadFilamentDataUnreadable = `24690096ff`
	loadFilamentData           = `{"transport":{"attrs":["request","twoway"],"id":9},"data":{"command":{"idx":52,"name":"load_initiate"}}};`
)

/*
func LoadFilament() {
	//Get all the stuff you need but it needs a different outpoint
	ctx, dev, _, done, iEp, oEp := generateUsbData()

	defer ctx.Close()
	defer dev.Close()
	defer done()

	log.Println("Loading Filament.....")

	cfg, _ := dev.Config(0)
	intf2, _ := cfg.Interface(0, 0)
	//Check active config num
	n, _ := dev.ActiveConfigNum()
	s := dev.ConfigDescription
	log.Println("Config Num", n)
	log.Println("Config Info", s)

	//You need a different endpoint for whatever reason so heres that
	oEp2, err := intf2.OutEndpoint(0x02)
	if err != nil {
		log.Fatalf("interface.OutEndpoint(0x02): %v", err)
	}

	//Generate correctly sized data
	data := []byte(loadFilamentDataUnreadable)
	// Write data to the USB device.

	//Unreadable stuff
	numBytes, _ := oEp2.Write(data)
	log.Printf("%d bytes sent", numBytes)

	//Readable Stuff
	data = []byte(loadFilamentData)
	numBytes, _ = oEp2.Write(data)
	log.Printf("%d bytes sent", numBytes)

	//Write data to start loading filament
	//dev.write(2, []byte(loadFilamentDataUnreadable))
	//dev.write(2, []byte(loadFilamentData))

	//Read the status forever
	//ReadStatusForever(iEp, oEp)
}
*/
