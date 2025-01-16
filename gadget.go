package gadget

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	GadgetConfigBasePath = "/sys/kernel/config/usb_gadget"
	UdcPathGlob          = "/sys/class/udc"
	StrEnglish           = "0x409"
)

// Gadget represents a USB gadget.
type Gadget struct {
	Name         string
	GadgetPath   string
	IdVendor     uint16
	IdProduct    uint16
	BcdDevice    uint16
	BcdUSB       uint16
	SerialNumber string
	Manufacturer string
	Product      string
	UDC          string
	Configs      []Config
}

// Config represents a USB gadget configuration.
type Config struct {
	Name          string
	Configuration string
	MaxPower      string
	Functions     []Function
}

// Function is an interface for USB gadget functions.
type Function interface {
	GadgetFunctionName() string
	GadgetFunctionCreate() Steps
}

// gadgetPath returns the path to the gadget.
func (g *Gadget) gadgetPath() string {
	if g.GadgetPath != "" {
		return g.GadgetPath
	}
	return filepath.Join(GadgetConfigBasePath, g.Name)
}

// Exists checks if the gadget exists.
func (g *Gadget) Exists() bool {
	_, err := os.Stat(g.gadgetPath())
	return !os.IsNotExist(err)
}

// Create creates the gadget.
func (g *Gadget) Create() error {
	return g.CreateSteps().Run()
}

// Remove removes the gadget.
func (g *Gadget) Remove() error {
	return g.RemoveSteps().Run()
}

// RemoveSteps generates steps to remove the gadget by reversing the creation steps.
func (g *Gadget) RemoveSteps() Steps {
	return g.CreateSteps().Undo().Reverse()
}

// CreateSteps generates steps to configure the gadget.
func (g *Gadget) CreateSteps() (steps Steps) {
	steps = Steps{
		Step{Mkdir, "", ""},
		Step{Write, "idVendor", fmt.Sprintf("0x%04x", g.IdVendor)},
		Step{Write, "idProduct", fmt.Sprintf("0x%04x", g.IdProduct)},
		Step{Write, "bcdUSB", fmt.Sprintf("0x%04x", g.BcdUSB)},
		Step{Write, "bcdDevice", fmt.Sprintf("0x%04x", g.BcdDevice)},

		Step{Mkdir, "strings/" + StrEnglish, ""},
		Step{Write, "strings/" + StrEnglish + "/serialnumber", g.SerialNumber},
		Step{Write, "strings/" + StrEnglish + "/manufacturer", g.Manufacturer},
		Step{Write, "strings/" + StrEnglish + "/product", g.Product},
	}

	for i := range g.Configs {
		var (
			c           = &g.Configs[i]
			configPath  = "configs/" + c.Name
			configSteps = Steps{
				Step{Comment, fmt.Sprintf("config `%s`", c.Name), ""},
				Step{Mkdir, "", ""},
				Step{Mkdir, "strings/" + StrEnglish, ""},
				Step{Write, "strings/" + StrEnglish + "/configuration", c.Configuration},
			}
		)

		configSteps.PrependPath(configPath)
		steps.Extend(configSteps)

		for _, fn := range c.Functions {
			var (
				name    = fn.GadgetFunctionName()
				fnPath  = "functions/" + name
				fnSteps = Steps{
					Step{Comment, fmt.Sprintf("config `%s`, function `%s`", c.Name, name), ""},
					Step{Mkdir, "", ""},
				}
			)

			fnSteps.Extend(fn.GadgetFunctionCreate())
			fnSteps.PrependPath(fnPath)
			steps.Extend(fnSteps)

			// Attach function to configuration
			steps.Append(Step{Symlink, fnPath, configPath + "/" + name})
		}
	}

	if v := g.UDC; v != "" {
		steps.Append(Step{Write, "UDC", v})
	}

	return steps.PrependPath(g.gadgetPath())
}
