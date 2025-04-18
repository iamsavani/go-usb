package gadget

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	Configs      map[string]*Config
}

// Config represents a USB gadget configuration.
type Config struct {
	Name          string
	Configuration string
	MaxPower      string
	Functions     map[string]Function
}

// Function is an interface for USB gadget functions.
type Function interface {
	GadgetFunctionName() string
	GadgetFunctionType() string
	GadgetFunctionCreate() Steps
}

func (g *Gadget) GetGadgetPath() string {
	if g.GadgetPath != "" {
		return g.GadgetPath
	}
	return filepath.Join(GadgetConfigBasePath, g.Name)
}

func (g *Gadget) ReadConfigfsFile(elem ...string) (string, error) {
	path := filepath.Join(g.GetGadgetPath(), filepath.Join(elem...))
	buf, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	return strings.TrimRight(string(buf), "\n"), nil
}

// Exists checks if the gadget exists.
func (g *Gadget) Exists() bool {
	_, err := os.Stat(g.GetGadgetPath())
	return !os.IsNotExist(err)
}

// Create creates the gadget.
func (g *Gadget) Create() error {
	return g.CreateSteps().Run()
}

// Remove removes the gadget.
func (g *Gadget) Remove() error {
	return g.CreateSteps().Reverse().Run()
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

	for name, c := range g.Configs {
		configPath := "configs/" + c.Name
		configSteps := Steps{
			Step{Comment, fmt.Sprintf("config `%s`", name), ""},
			Step{Mkdir, "", ""},
			Step{Mkdir, "strings/" + StrEnglish, ""},
			Step{Write, "strings/" + StrEnglish + "/configuration", c.Configuration},
		}
		configSteps.PrependPath(configPath)
		steps.Extend(configSteps)

		for fname, fn := range c.Functions {
			fnPath := "functions/" + fname
			fnSteps := Steps{
				Step{Comment, fmt.Sprintf("config `%s`, function `%s`", name, fname), ""},
				Step{Mkdir, "", ""},
			}
			fnSteps.Extend(fn.GadgetFunctionCreate())
			fnSteps.PrependPath(fnPath)
			steps.Extend(fnSteps)
			// Attach function to configuration
			steps.Append(Step{Symlink, fnPath, configPath + "/" + fname})
		}
	}

	if g.UDC != "" {
		steps.Append(Step{Write, "UDC", g.UDC})
	}

	return steps.PrependPath(g.GetGadgetPath())
}

func (g *Gadget) GetFunctionPath(fnName string) (string, bool) {
	for _, cfg := range g.Configs {
		if _, ok := cfg.Functions[fnName]; ok {
			return filepath.Join(g.GetGadgetPath(), "functions", fnName), true
		}
	}
	return "", false
}

// Unbind disables the gadget by writing an empty string to the UDC file.
func (g *Gadget) Unbind() error {
	udcPath := filepath.Join(g.GetGadgetPath(), "UDC")
	err := os.WriteFile(udcPath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("failed to unbind UDC: %w", err)
	}
	return nil
}

// Bind enables the gadget by writing the UDC name to the UDC file.
func (g *Gadget) Bind(udc string) error {
	if udc != "" {
		g.UDC = udc
	} else {
		udcs := GetUdcs()
		if len(udcs) < 1 {
			return fmt.Errorf("no UDC found")
		}
		g.UDC = udcs[0]
	}

	udcPath := filepath.Join(g.GetGadgetPath(), "UDC")
	if err := os.WriteFile(udcPath, []byte(g.UDC), 0644); err != nil {
		return fmt.Errorf("failed to bind UDC: %w", err)
	}

	return nil
}
