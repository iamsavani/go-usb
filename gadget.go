package gadget

import (
	"fmt"
	"io/ioutil"
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
func (g *Gadget) GetGadgetPath() string {
	if g.GadgetPath != "" {
		return g.GadgetPath
	}
	return filepath.Join(GadgetConfigBasePath, g.Name)
}

func (g *Gadget) ReadConfigfsFile(elem ...string) (string, error) {
	var path = filepath.Join(g.GetGadgetPath(), filepath.Join(elem...))
	if buf, err := ioutil.ReadFile(path); err != nil {
		return "", nil
	} else {
		return strings.TrimRight(string(buf), "\n"), nil
	}
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
	return g.BuildRemovalSteps().Run()
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

	return steps.PrependPath(g.GetGadgetPath())
}

// AddFunction dynamically adds a function to the first config of the gadget
func (g *Gadget) AddFunction(functionType, additionalConfig interface{}) error {
	if len(g.Configs) == 0 {
		return fmt.Errorf("no configuration found to add function")
	}

	var newFunction Function

	// Choose function type based on the passed type
	switch functionType {
	case "hid":
		// Assuming additionalConfig contains the type of HID function (keyboard or mouse)
		if config, ok := additionalConfig.(HidFunction); ok {
			newFunction = &config
		} else {
			return fmt.Errorf("invalid configuration for HID function")
		}
	case "mass_storage":
		// Assuming additionalConfig contains the specific MassStorageFunction details
		if config, ok := additionalConfig.(MassStorageFunction); ok {
			newFunction = &config
		} else {
			return fmt.Errorf("invalid configuration for mass storage")
		}
	default:
		return fmt.Errorf("unsupported function type: %s", functionType)
	}

	// Add the function to the first config's function list
	g.Configs[0].Functions = append(g.Configs[0].Functions, newFunction)

	// Regenerate gadget steps to apply the changes
	return g.CreateSteps().Run()
}

func (g *Gadget) RemoveFunction(fnName string) error {
	for _, cfg := range g.Configs {
		configPath := filepath.Join(g.GetGadgetPath(), "configs", cfg.Name)
		symlinkPath := filepath.Join(configPath, fnName)

		if _, err := os.Lstat(symlinkPath); err == nil {
			if err := os.Remove(symlinkPath); err != nil {
				return err
			}
			fnPath := filepath.Join(g.GetGadgetPath(), "functions", fnName)
			if err := os.RemoveAll(fnPath); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("function %s not found in any config", fnName)
}

// RemoveAllFunctions removes all functions from all configs
func (g *Gadget) RemoveAllFunctions() error {
	for _, cfg := range g.Configs {
		for _, fn := range cfg.Functions {
			if err := g.RemoveFunction(fn.GadgetFunctionName()); err != nil {
				return err
			}
		}
	}
	return nil
}

// BuildRemovalSteps constructs Steps to fully remove the gadget from configfs
func (g *Gadget) BuildRemovalSteps() Steps {
	var steps Steps
	gadgetPath := g.GetGadgetPath()

	// Unbind
	steps.Append(Step{Write, "UDC", ""})

	for _, cfg := range g.Configs {
		configPath := filepath.Join("configs", cfg.Name)

		for _, fn := range cfg.Functions {
			fnName := fn.GadgetFunctionName()
			fnPath := filepath.Join("functions", fnName)

			steps.Append(Step{Remove, filepath.Join(configPath, fnName), ""})
			steps.Append(Step{Remove, fnPath, ""})
		}

		steps.Append(Step{Remove, filepath.Join(configPath, "strings/"+StrEnglish), ""})
		steps.Append(Step{Remove, configPath, ""})
	}

	steps.Append(Step{Remove, filepath.Join("strings", StrEnglish), ""})
	steps.Append(Step{Remove, "", ""}) // remove root gadget dir

	return steps.PrependPath(gadgetPath)
}

func (g *Gadget) GetFunctionPath(functionName string) (string, bool) {
	for _, cfg := range g.Configs {
		for _, fn := range cfg.Functions {
			if fn.GadgetFunctionName() == functionName {
				return filepath.Join(g.GetGadgetPath(), "functions", functionName), true
			}
		}
	}
	return "", false
}
