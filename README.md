# go-usb

`go-usb` is a Go package for creating and managing USB gadget configurations on Linux systems. This package allows you to define and configure various USB gadget functions such as HID and Mass Storage.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [HID Function](#hid-function)
  - [Mass Storage Function](#mass-storage-function)
- [Example](#example)
- [Acknowledgements](#acknowledgements)
- [License](#license)

## Installation

To install the `go-usb` package, use the following command:

```sh
go get github.com/iamsavani/go-usb
```

## Usage

### HID Function

The `HidFunction` struct represents a HID gadget function. You can create a new HID function and generate the steps required to configure it.


#### Example
```go
package main

import (
    "github.com/iamsavani/go-usb"
    "fmt"
)

func main() {
    hid := usb.HidFunction{
        Name:         "example",
        Protocol:     1,
        Subclass:     1,
        ReportLength: 8,
        Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
    }

    steps := hid.GadgetFunctionCreate()
    fmt.Println("HID Function Steps:", steps)
}
```

### Mass Storage Function


The `MassStorageFunction` struct represents a Mass Storage gadget function. You can create a new Mass Storage function with one or more logical unit numbers (LUNs).

#### Example
```go
package main

import (
    "github.com/iamsavani/go-usb"
    "fmt"
)

func main() {
    lun := usb.MassStorageLun{
        Name:          "example_lun",
        File:          "/path/to/file",
        Removable:     true,
        Cdrom:         false,
        Ro:            false,
        InquiryString: "Example Inquiry",
    }

    ms := usb.MassStorageFunction{
        Name:  "example",
        Stall: true,
        Luns:  []usb.MassStorageLun{lun},
    }

    steps := ms.GadgetFunctionCreate()
    fmt.Println("Mass Storage Function Steps:", steps)
}
```

## Example

An example usage of the `go-usb` package is provided in the `example` folder.


Here is a complete example that combines both HID and Mass Storage functions in a single gadget:

```go
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
```

## Acknowledgements

This project was inspired by and uses code from the [go-linuxuapi](https://github.com/pdmccormick/go-linuxuapi) repository by pdmccormick. Special thanks to the original author for their work.


## License

This project is licensed under the MIT License - see the [LICENSE](/LICENSE) file for details.

