package modt

import (
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/google/gousb"
	"log"
	"time"
)

const (
	MTVID = 0x2b75
	MTPID = 0x0002

	STATUS_STRING = `{"metadata":{"version":1,"type":"status"}}`
)

type Comm int

const (
	COMM_PAUSE Comm = iota
	COMM_PLAY
	COMM_END
)

var (
	ctx      *gousb.Context
	dev      *gousb.Device
	USBInt   *gousb.Interface
	doneFunc func()
	iEp      *gousb.InEndpoint
	oEp      *gousb.OutEndpoint
)

func SaveUsbData() {
	var err error
	ctx, dev, USBInt, doneFunc, iEp, oEp, err = generateUsbData()
	if err != nil {
		log.Println(err)
		log.Println("Device not found. Probably Disconnected")
		return
	}
}

func GetStatus(controlChan *chan Comm, info *PrinterInformation) string {
	defer ctx.Close()
	defer dev.Close()
	defer doneFunc()
	defer func() { fmt.Println("FINISHED DEFERRING") }()

	//Read the status forever
	ReadStatusForever(iEp, oEp, controlChan, info)
	return "Ended connection? (This really shouldnt happen)"
}

func ReadStatusForever(iEp *gousb.InEndpoint, oEp *gousb.OutEndpoint, commChan *chan Comm, info *PrinterInformation) {
	doing := false
	for {
		//Check if checking should happen
		select {
		case do := <-*commChan:
			//doing = do
			if do == COMM_PLAY {
				doing = true
				fmt.Println("Will Log")
			} else if do == COMM_PAUSE {
				doing = false
				fmt.Println("Will Not Log")
			} else if do == COMM_END { //End transmission and hopefully things will defer
				fmt.Println("Ending transmission")
				goto EndTransmissionLabel
			}
		default:
			if doing {
				fmt.Println("Asking....")
				//Ask for data
				//Generate correctly sized data
				data := []byte(STATUS_STRING) //make([]byte,0)

				// Write data to the USB device.
				numBytes, err := oEp.Write(data)
				if err != nil {
					log.Fatalf("Problem sending, ", err)
				}
				log.Printf("%d bytes sent", numBytes)

				fmt.Println("Finished asking...Listening...")

				//Read, and also send down a channel
				text := readModt(iEp)

				//var info PrinterInformation
				json.Unmarshal([]byte(text), &info)
				g.Update()
				fmt.Println(info)

				//Wait
				time.Sleep(1 * time.Second)
			}
		}
	}
EndTransmissionLabel:
}

func readModt(epIn *gousb.InEndpoint) string {
	//Ask and you shall recieve
	buf := make([]byte, 64)
	fmt.Println("Made buffer")

	readBytes, err := epIn.Read(buf)
	if err != nil {
		log.Fatalf("Error Reading: %v", err)
	}
	if readBytes == 0 {
		log.Fatalf("In Endpoint returned 0 bytes")
	}

	text := string(buf)
	fulltext := ""
	fulltext += text
	log.Println("Read once")

	for readBytes == 64 {
		readBytes, err = epIn.Read(buf)
		if err != nil {
			log.Fatalf("Error Reading: %v", err)
		}
		if readBytes == 0 {
			log.Fatalf("In Endpoint returned 0 bytes")
		}
		text = string(buf[0:readBytes])
		fulltext += text
	}
	return fulltext
}

//generateUsbData gets all usb data needed
func generateUsbData() (*gousb.Context, *gousb.Device, *gousb.Interface, func(), *gousb.InEndpoint, *gousb.OutEndpoint, error) {
	// Initialize a new Context.
	ctx := gousb.NewContext()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(MTVID, MTPID)
	if err != nil {
		log.Println("PROBLEM")
		log.Fatalf("Could not open a device: %v", err)
	}
	if dev == nil {
		fmt.Println("This is a problem")
		log.Println(dev)
		return nil, nil, nil, nil, nil, nil, errors.New("No Device Found")
	}

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active

	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	}

	//Open an in endpoint
	iEp, err := intf.InEndpoint(0x83)
	if err != nil {
		log.Fatalf("%s.InEndpoint(0x83): %v", intf, err)
	}

	// Open an OUT endpoint.
	oEp, err := intf.OutEndpoint(0x04)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(0x04): %v", intf, err)
	}
	return ctx, dev, intf, done, iEp, oEp, nil

}
