package main

import (
	"fmt"
	"log"

	gadget "github.com/iamsavani/gadget"
)

func main() {
	kbd := gadget.HidFunction{
		Name:         "usb0",
		Protocol:     1,
		Subclass:     0,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}

	mouse := gadget.HidFunction{
		Name:         "usb1",
		Protocol:     2,
		Subclass:     0,
		ReportLength: 6,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}

	mass := gadget.MassStorageFunction{
		Name: "usb0",
		Luns: []gadget.MassStorageLun{
			gadget.MassStorageLun{
				Name:          "0",
				File:          "\n",
				Removable:     true,
				Cdrom:         true,
				Ro:            true,
				InquiryString: "Mass Storage",
			},
		},
	}

	// Create a Gadget
	g := gadget.Gadget{
		Name:         "gadget",
		IdVendor:     0x1d6b,
		IdProduct:    0x0104,
		BcdDevice:    0x0100,
		BcdUSB:       0x0200,
		SerialNumber: "1234567890",
		Manufacturer: "Example Manufacturer",
		Product:      "Example Product",
		Configs: []gadget.Config{
			{
				Name:          "c.1",
				Configuration: "Config 1: HID",
				MaxPower:      "250",
				Functions:     []gadget.Functions{&kbd, &mouse, &mass},
			},
		},
	}

	if err := g.Create(); err != nil {
		log.Fatalf("error creating: %s", err)
	}

	fmt.Println("created")
}
