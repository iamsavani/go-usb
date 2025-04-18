package gadget

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFunction(t *testing.T) {

	// Create a mass storage function
	massFunc := &MassStorageFunction{
		Name:    "mass_storage.usb0",
		Stall:   true,
		Enabled: true,
		Luns: []MassStorageLun{
			{
				Name: "0",
				Attr: GadgetAttributes{
					"file":           "\n",
					"removable":      "1",
					"cdrom":          "1",
					"ro":             "1",
					"inquiry_string": "Test Mass Storage",
				},
			},
		},
	}

	// Setup a test gadget
	gadget := &Gadget{
		Name: "test_gadget_add",
		Attrs: GadgetAttributes{
			"idVendor":  "0x1234",
			"idProduct": " 0x5678",
		},
		ConfigAttrs: GadgetAttributes{
			"serial_number": "123456789",
			"manufacturer":  "TestManufacturer",
			"product":       "TestProduct",
		},
		Configs: map[string]*Config{
			"c.1": {
				Name:          "c.1",
				Configuration: "Mass Storage Config",
				MaxPower:      "120",
				Functions: map[string]Function{
					"mass_storage.usb0": massFunc,
				},
			},
		},
	}

	// Create the gadget (if not already created)
	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		// return
	}

	// Create the gadget
	err := gadget.Create()
	assert.NoError(t, err, "Error creating gadget")

	// Verify the symlink was created correctly
	fnPath := filepath.Join(gadget.GetGadgetPath(), "functions", massFunc.Name)
	configPath := filepath.Join(gadget.GetGadgetPath(), "configs", "c.1")
	_, err = os.Stat(fnPath)
	assert.NoError(t, err, fmt.Sprintf("Function path %s does not exist", fnPath))

	_, err = os.Stat(configPath)
	assert.NoError(t, err, fmt.Sprintf("Config path %s does not exist", configPath))

	// Check if the symlink exists inside the config directory
	symlinkTarget := filepath.Join(configPath, massFunc.Name)
	symlink, err := os.Readlink(symlinkTarget)
	assert.NoError(t, err, fmt.Sprintf("Symlink %s does not exist", symlinkTarget))
	assert.Equal(t, fnPath, symlink, fmt.Sprintf("Symlink target mismatch for %s", symlinkTarget))

	// Cleanup after the test
	err = gadget.Remove()
	assert.NoError(t, err, "Error removing gadget")
}
