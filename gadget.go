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
	buf, err := ioutil.ReadFile(path)
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

// AddFunction adds a new function to the gadget.
func (g *Gadget) AddFunction(configName, functionName string, fn Function) error {
	cfg, ok := g.Configs[configName]
	if !ok {
		return fmt.Errorf("config %s not found", configName)
	}
	if _, exists := cfg.Functions[functionName]; exists {
		return fmt.Errorf("function %s already exists", functionName)
	}
	cfg.Functions[functionName] = fn

	fnPath := "functions/" + functionName
	configPath := "configs/" + configName

	steps := Steps{
		Step{Comment, fmt.Sprintf("Add function `%s` to config `%s`", functionName, configName), ""},
		Step{Mkdir, "", ""},
	}
	steps.Extend(fn.GadgetFunctionCreate())
	steps.PrependPath(fnPath)

	steps.Append(Step{Symlink, fnPath, configPath + "/" + functionName})

	steps.PrependPath(g.GetGadgetPath())
	return steps.Run()
}

// BuildRemovalSteps returns steps to remove the entire gadget configuration.
func (g *Gadget) BuildRemovalSteps() Steps {
	steps := Steps{Step{Write, "UDC", ""}}
	for cfgName, cfg := range g.Configs {
		configPath := filepath.Join("configs", cfgName)
		for fnName := range cfg.Functions {
			functionPath := filepath.Join("functions", fnName)
			steps.Append(Step{Remove, filepath.Join(configPath, fnName), ""})
			steps.Append(Step{Remove, functionPath, ""})
		}
		steps.Append(Step{Remove, filepath.Join(configPath, "strings", StrEnglish), ""})
		steps.Append(Step{Remove, configPath, ""})
	}
	steps.Append(Step{Remove, filepath.Join("strings", StrEnglish), ""})
	steps.Append(Step{Remove, "", ""}) // remove root gadget dir
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
