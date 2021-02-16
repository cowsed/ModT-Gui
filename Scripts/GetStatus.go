package modt

import (
	"fmt"
	"github.com/google/gousb"
	"log"
	"time"
	"errors"
)

const (
	MTVID = 0x2b75
	MTPID = 0x0002

	STATUS_STRING = `{"metadata":{"version":1,"type":"status"}}`
)

func GetStatus(controlChan *chan bool) string {
	ctx, dev, _, done, iEp, oEp, err := generateUsbData()
	if err!=nil{
		log.Println(err)
		return "Device not found. Probably Disconnected"
	}
	defer ctx.Close()
	defer dev.Close()
	defer done()

	//Read the status forever
	ReadStatusForever(iEp, oEp, controlChan)
	return "Ended connection? (This really shouldnt happen)"
}

func ReadStatusForever(iEp *gousb.InEndpoint, oEp *gousb.OutEndpoint, commChan *chan bool) {
	doing := false
	for {
		//Check if checking should happen
		select {
		case do := <-*commChan:
			doing = do
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
				fmt.Println(text)
				//Wait
				time.Sleep(2 * time.Second)
			}
		}
	}

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
		return nil,nil,nil,nil,nil,nil,errors.New("No Device Found")
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
	return ctx, dev, intf, done, iEp, oEp,nil

}
