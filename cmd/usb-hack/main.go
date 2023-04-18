package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
)

func main() {
	// Initialize a new Context.
	ctx := gousb.NewContext()
	// Ensure that the Context is closed before the program exits.
	// defer ctx.Close()
	
	// Open a device with a specific VID/PID.
	dev, err := ctx.OpenDeviceWithVIDPID(0x046d, 0xc21d)
	if err != nil {
		log.Fatalf("Could not open device: %v", err)
	}
	fmt.Println("Device opened: ", usbid.Describe(dev.Desc))

	// Ensure that the device is closed before the program exits.
	// defer dev.Close()
	dev.SetAutoDetach(true)

	// Claim the default interface.
	// The default interface is always #0 alt #0 in the currently active config
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", usbid.Describe(dev.Desc), err)
	}
	fmt.Println("Interface claimed: ", intf)
	// Ensure that the interface is released before the program exits.
	// defer done()

	// Open an IN endpoint.
	// The endpoint number is always 0x81 for the first IN endpoint.
	ep, err := intf.InEndpoint(0x81)
	if err != nil {
		log.Fatalf("%s.InEndpoint(0x81): %v", intf, err)
	}
	fmt.Println("Endpoint opened: ", ep)

	// Read from the endpoint.
	bufsize := intf.Setting.Endpoints[0x81].MaxPacketSize
	buf := make([]byte, bufsize)
	go readEp(ep, buf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	fmt.Println("Got signal:", sig)
	fmt.Print("Releasing interface...")
	done()
	fmt.Println("done")
	fmt.Print("Closing device...")
	dev.Close()
	fmt.Println("closed")
	fmt.Print("Closing context...")
	ctx.Close()
	fmt.Println("closed")

	return
}

func readEp(ep *gousb.InEndpoint, buf []byte) {
	for {
		n, err := ep.Read(buf)
		if err != nil {
			log.Fatalf("%s.Read(): %v", ep, err)
		}
		fmt.Printf("Read %d bytes: %v\n", n, buf)
	}
}
