package modt

import (
	"fmt"
)

const MOD_ADLER uint32 = 65521

func SendGcode(file string) {
	checksum := adler32(file)
	fmt.Println("Checksum:", checksum)

	// These came from usb dump.
	// Some commands are human readable some are maybe checksums

	/*
		fmt.Println("Configs:",dev.Desc.Configs)

		n,_:=dev.ActiveConfigNum()
		fmt.Println("Active Config:",n)
		//dev
		n,_=dev.ActiveConfigNum()
		fmt.Println("Active Config:",n)

	*/

	numBytes, err := oEp2.Write([]byte{0x24, 0x6a, 0x00, 0x95, 0xff})
	checkSend(err)

	numBytes, err = oEp2.Write([]byte(`{"transport":{"attrs":["request","twoway"],"id":3},"data":{"command":{"idx":0,"name":"bio_get_version"}}};`))
	checkSend(err)

	fmt.Println("Reading from 81")

	//For whatever reason this just blocks here but works fine in the python version
	fmt.Println(readModt(iEp81))
	fmt.Println("Finished Reading from 81")

	//Get Normal Status
	numBytes, err = oEp4.Write([]byte(STATUS_STRING))
	fmt.Println(readModt(iEp83))

	numBytes, err = oEp2.Write([]byte{0x24, 0x8b, 0x00, 0x74, 0xff})
	checkSend(err)
	numBytes, err = oEp2.Write([]byte(`{"transport":{"attrs":["request","twoway"],"id":5},"data":{"command":{"idx":22,"name":"wifi_client_get_status","args":{"interface_t":0}}}};`))
	checkSend(err)

	fmt.Println("Reading from 81 - 2")
	fmt.Println(readModt(iEp81))
	fmt.Println("Finished Reading from 81 - 2")

	numBytes, err = oEp2.Write([]byte{0x24, 0x8b, 0x00, 0x74, 0xff})
	checkSend(err)
	numBytes, err = oEp2.Write([]byte(`{"transport":{"attrs":["request","twoway"],"id":5},"data":{"command":{"idx":22,"name":"wifi_client_get_status","args":{"interface_t":0}}}};`))
	checkSend(err)

	fmt.Println("Reading from 81 - 3")
	fmt.Println(readModt(iEp81))
	fmt.Println("Finished Reading from 81 - 3")

	fmt.Println(numBytes)
	//Get Normal Status
	numBytes, err = oEp4.Write([]byte(STATUS_STRING))
	checkSend(err)
	fmt.Println(readModt(iEp83))

	//Get Normal Status
	numBytes, err = oEp4.Write([]byte(STATUS_STRING))
	checkSend(err)
	fmt.Println(readModt(iEp83))

	// Start writing actual gcode
	// File size and adler32 checksum calculated earlier
	size := len(file)

	command := fmt.Sprintf(`{"metadata":{"version":1,"type":"file_push"},"file_push":{"size":'%d',"adler32":'%d',"job_id":""}}`, size, checksum)
	fmt.Println("Commanding: ", command)
	numBytes, err = oEp4.Write([]byte(command))
	fmt.Println("command num bytes:",numBytes)
	if err!=nil{
		fmt.Println("Command Err",err)
	}

	// Write gcode in batches of 20 bulk writes, each 5120 bytes.
	// Read mod-t status between these 20 bulk writes

	/*
			start=0
		counter=0
		while True:
		 if (start+5120-1>size-1):
		        end=size
		 else:
		        end=start+5120
		 block = gcode[start:end]
		 print(str(counter)+':' +str(start)+'-'+str(end-1)+'\t'+str(len(block)))
		 counter += 1
		 if counter>=20:
		  temp=read_modt(0x83)
		  counter = 0
		 dev.write(4, block)
		 if (start == 0):
		  temp=read_modt(0x83)
		 start = start + 5120
		 if (start>size):
		        break;

	*/
	fmt.Println("GCODE Size", size)
	start := 0
	counter := 0
	end := 0
	var temp string
	for {
		fmt.Println("Sending")
		if start+5120 > size-1 {
			end = size
		} else {
			end = start + 5120
		}
		block := []byte(file[start:end])
		fmt.Printf("%d: %d-%d \t%d\n", counter, start, end-1, len(block))
		counter++
		if counter >= 20 {
			fmt.Println("querying statuss")
			temp = readModt(iEp83)
			counter = 0
		}
		fmt.Println("Writing Block")
		numBytes, err := oEp4.Write(block)
		fmt.Printf("Sent %d of %d bytes\n", numBytes, len(block))
		checkSend(err)

		if start == 0 {
			fmt.Println("Beginning Read")
			temp = readModt(iEp83)
			fmt.Println(temp)
		}
		start = start + 5120
		if start > size {
			break
		}
		fmt.Println(temp)
	}

}
func checkSend(err error) {
	if err != nil {
		fmt.Println("Problem sending", err)
	}
}

//Generate Adler32 Checksum of file
func adler32(data string) uint32 {
	var a uint32 = 0
	var b uint32 = 0

	// Process each byte of the data in order
	// Bad because its not necessarily ascii in go but whatever
	for _, d := range []byte(data) {
		a = (a + uint32(d)) % MOD_ADLER
		b = (b + a) % MOD_ADLER
	}
	return (b << 16) | a
}
