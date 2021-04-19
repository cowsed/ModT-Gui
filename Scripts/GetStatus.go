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
	conf     *gousb.Config
	doneFunc func()
	iEp83    *gousb.InEndpoint
	iEp81    *gousb.InEndpoint
	oEp4     *gousb.OutEndpoint
	oEp2     *gousb.OutEndpoint
)

func SaveUsbData() {
	var err error
	ctx, dev, conf, USBInt, doneFunc, iEp83, iEp81, oEp4, oEp2, err = generateUsbData()
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
	defer conf.Close()
	defer func() { fmt.Println("FINISHED DEFERRING") }()

	//Read the status forever
	ReadStatusForever(iEp83, oEp4, controlChan, info)
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
				fmt.Println("Am here for a good? reason")
				//Ask for data
				//Generate correctly sized data
				data := []byte(STATUS_STRING) //make([]byte,0)

				// Write data to the USB device.
				_, err := oEp.Write(data)
				if err != nil {
					log.Fatalf("Problem sending, %v", err)
				}
				//log.Printf("%d bytes sent", numBytes)

				//Read, and also send down a channel
				text := readModt(iEp83)

				//var info PrinterInformation
				json.Unmarshal([]byte(text), &info)
				g.Update()

				//Wait
				time.Sleep(1 * time.Second)
			}
		}
	}
EndTransmissionLabel:
}

func readModt(in *gousb.InEndpoint) string {
	//Ask and you shall recieve
	buf := make([]byte, 64)

	readBytes, err := in.Read(buf)
	if err != nil {
		log.Fatalf("Error Reading: %v", err)
	}
	if readBytes == 0 {
		log.Fatalf("In Endpoint returned 0 bytes")
	}

	text := string(buf)
	fulltext := ""
	fulltext += text

	for readBytes == 64 {
		readBytes, err = in.Read(buf)
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
func generateUsbData() (*gousb.Context, *gousb.Device, *gousb.Config, *gousb.Interface, func(), *gousb.InEndpoint, *gousb.InEndpoint, *gousb.OutEndpoint, *gousb.OutEndpoint, error) {
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
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, errors.New("No Device Found")
	}

	// Claim the default interface using a convenience function.
	//..// The default interface is always #0 alt #0 in the currently active
	conf, err := dev.Config(1)
	if err != nil {
		panic(err)
	}

	intf, err := conf.Interface(0, 0)
	fmt.Println("DESCRIPTION:", conf.String())
	//conf := nil
	//intf, done, err := dev.DefaultInterface()

	if err != nil {
		log.Fatalf("%s getting interface failed: %v", dev, err)
	}

	//Open an in endpoint for regular reading
	iEp83, err := intf.InEndpoint(0x83)
	if err != nil {
		log.Fatalf("%s.InEndpoint(0x83): %v", intf, err)
	}
	iEp81, err := intf.InEndpoint(0x81)
	if err != nil {
		log.Fatalf("%s.InEndpoint(0x81): %v", intf, err)
	}

	// Open an OUT endpoint. for statusing?
	oEp4, err := intf.OutEndpoint(0x04)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(0x04): %v", intf, err)
	}
	// Open an OUT endpoint. for other things?
	oEp2, err := intf.OutEndpoint(0x02)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(0x02): %v", intf, err)
	}

	done := func() {}
	return ctx, dev, conf, intf, done, iEp83, iEp81, oEp4, oEp2, nil

}
