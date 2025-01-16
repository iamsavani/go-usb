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

// CreateSteps generates steps to configure the gadget.
func (g *Gadget) CreateSteps() (steps Steps) {
    steps = Steps{
        {Action: Mkdir, Path: "", Value: ""},
        {Action: Write, Path: "idVendor", Value: fmt.Sprintf("0x%04x", g.IdVendor)},
        {Action: Write, Path: "idProduct", Value: fmt.Sprintf("0x%04x", g.IdProduct)},
        {Action: Write, Path: "bcdUSB", Value: fmt.Sprintf("0x%04x", g.BcdUSB)},
        {Action: Write, Path: "bcdDevice", Value: fmt.Sprintf("0x%04x", g.BcdDevice)},
        {Action: Mkdir, Path: "strings/" + StrEnglish, Value: ""},
        {Action: Write, Path: "strings/" + StrEnglish + "/serialnumber", Value: g.SerialNumber},
        {Action: Write, Path: "strings/" + StrEnglish + "/manufacturer", Value: g.Manufacturer},
        {Action: Write, Path: "strings/" + StrEnglish + "/product", Value: g.Product},
    }

    for i := range g.Configs {
        c := &g.Configs[i]
        configPath := "configs/" + c.Name
        configSteps := Steps{
            {Action: Comment, Path: fmt.Sprintf("config `%s`", c.Name), Value: ""},
            {Action: Mkdir, Path: "", Value: ""},
            {Action: Mkdir, Path: "strings/" + StrEnglish, Value: ""},
            {Action: Write, Path: "strings/" + StrEnglish + "/configuration", Value: c.Configuration},
        }

        configSteps.PrependPath(configPath)
        steps.AddSteps(configSteps)

        for _, fn := range c.Functions {
            name := fn.GadgetFunctionName()
            fnPath := "functions/" + name
            fnSteps := Steps{
                {Action: Comment, Path: fmt.Sprintf("config `%s`, function `%s`", c.Name, name), Value: ""},
                {Action: Mkdir, Path: "", Value: ""},
            }

            fnSteps.AddSteps(fn.GadgetFunctionCreate())
            fnSteps.PrependPath(fnPath)
            steps.AddSteps(fnSteps)

            // Attach function to configuration
            steps.Add(Step{Action: Symlink, Path: fnPath, Value: configPath + "/" + name})
        }
    }

    if g.UDC != "" {
        steps.Add(Step{Action: Write, Path: "UDC", Value: g.UDC})
    }

    return steps.PrependPath(g.gadgetPath())
}