package modt

import (
	"fmt"
)

var UNREADABLE = [...]byte{0x24, 0x6a, 0x00, 0x95, 0xff}
var SEND_STRING = `{"transport":{"attrs":["request","twoway"],"id":3},"data":{"command":{"idx":0,"name":"bio_get_version"}}}`

const MOD_ADLER uint32 = 65521

func SendGcode(file string) {
	checksum := adler32(file)
	fmt.Println("Checksum:", checksum)
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
