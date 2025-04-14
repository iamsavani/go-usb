package gadget

import (
	"fmt"
	"testing"
)

// TestAddFunction tests the dynamic addition of HID and mass storage functions to a gadget.
func TestAddFunction(t *testing.T) {
	// Create a new Gadget instance
	gadget := &Gadget{
		Name:         "test_gadget_add",
		IdVendor:     0x1234,
		IdProduct:    0x5678,
		SerialNumber: "123456789",
		Manufacturer: "TestManufacturer",
		Product:      "TestProduct",
	}

	// Create a configuration and add it to the gadget
	config := Config{
		Name:          "config1",
		Configuration: "Test Configuration",
		MaxPower:      "500",
		Functions:     []Function{},
	}
	gadget.Configs = append(gadget.Configs, config)

	// Test adding HID keyboard function
	keyboardHid := HidFunction{
		Name:         "keyboard",
		Protocol:     0x01,
		Subclass:     0x01,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}

	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		return
	}

	err := gadget.AddFunction("hid", keyboardHid)
	if err != nil {
		t.Errorf("Error adding HID keyboard function: %v", err)
	}

	// Check if the HID function is correctly added
	fnPath, exists := gadget.GetFunctionPath("hid.keyboard")
	if !exists {
		t.Errorf("Expected function path for HID keyboard, but it does not exist.")
	} else {
		t.Logf("HID keyboard function path: %s", fnPath)
	}

	// Test adding HID mouse function
	mouseHid := HidFunction{
		Name:         "mouse",
		Protocol:     0x02,
		Subclass:     0x02,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x02},
	}
	err = gadget.AddFunction("hid", mouseHid)
	if err != nil {
		t.Errorf("Error adding HID mouse function: %v", err)
	}

	// Check if the HID mouse function is correctly added
	fnPath, exists = gadget.GetFunctionPath("hid.mouse")
	if !exists {
		t.Errorf("Expected function path for HID mouse, but it does not exist.")
	} else {
		t.Logf("HID mouse function path: %s", fnPath)
	}

	// Test adding Mass Storage function
	massFunc := MassStorageFunction{
		Name:  "usb0",
		Stall: true,
		Luns: []MassStorageLun{
			{
				Name:          "0",
				File:          "\n", // Replace with actual device path
				Removable:     true,
				Cdrom:         true,
				Ro:            true,
				InquiryString: "Test Mass Storage",
			},
		},
	}
	err = gadget.AddFunction("mass_storage", massFunc)
	if err != nil {
		t.Errorf("Error adding mass storage function: %v", err)
	}

	// Check if the Mass Storage function is correctly added
	fnPath, exists = gadget.GetFunctionPath("mass_storage.usb0")
	if !exists {
		t.Errorf("Expected function path for Mass Storage usb0, but it does not exist.")
	} else {
		t.Logf("Mass Storage function path: %s", fnPath)
	}
}

// TestInvalidAddFunction tests invalid function types passed to AddFunction.
func TestInvalidAddFunction(t *testing.T) {
	// Create a new Gadget instance
	gadget := &Gadget{
		Name:         "test_gadget_invalid",
		IdVendor:     0x1234,
		IdProduct:    0x5678,
		SerialNumber: "123456789",
		Manufacturer: "TestManufacturer",
		Product:      "TestProduct",
	}

	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		return
	}

	// Test adding an invalid function type
	err := gadget.AddFunction("invalid_type", nil)
	if err == nil {
		t.Errorf("Expected error when adding function with invalid type, but got none.")
	} else {
		t.Logf("Correctly caught error: %v", err)
	}
}

// TestFunctionNameGeneration tests the function name generation logic.
func TestFunctionNameGeneration(t *testing.T) {
	// Create a new Gadget instance
	gadget := &Gadget{
		Name:         "test_gadget_fn_name",
		IdVendor:     0x1234,
		IdProduct:    0x5678,
		SerialNumber: "123456789",
		Manufacturer: "TestManufacturer",
		Product:      "TestProduct",
	}

	// Create a configuration and add it to the gadget
	config := Config{
		Name:          "config1",
		Configuration: "Test Configuration",
		MaxPower:      "500",
		Functions:     []Function{},
	}
	gadget.Configs = append(gadget.Configs, config)

	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		return
	}

	// Add HID keyboard function
	keyboardHid := HidFunction{
		Name:         "keyboard",
		Protocol:     0x01,
		Subclass:     0x01,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}
	err := gadget.AddFunction("hid", keyboardHid)
	if err != nil {
		t.Errorf("Error adding HID keyboard function: %v", err)
	}

	// Test if the correct function name is generated for HID keyboard
	fnPath, exists := gadget.GetFunctionPath("hid.keyboard")
	if !exists {
		t.Errorf("Expected function path for HID keyboard, but it does not exist.")
	} else {
		t.Logf("HID keyboard function path: %s", fnPath)
	}

	// Add Mass Storage function
	massFunc := MassStorageFunction{
		Name:  "usb0",
		Stall: true,
		Luns: []MassStorageLun{
			{
				Name:          "0",
				File:          "\n",
				Removable:     true,
				Cdrom:         true,
				Ro:            true,
				InquiryString: "Test Mass Storage",
			},
		},
	}
	err = gadget.AddFunction("mass_storage", massFunc)
	if err != nil {
		t.Errorf("Error adding mass storage function: %v", err)
	}

	// Test if the correct function name is generated for Mass Storage
	fnPath, exists = gadget.GetFunctionPath("mass_storage.usb0")
	if !exists {
		t.Errorf("Expected function path for Mass Storage usb0, but it does not exist.")
	} else {
		t.Logf("Mass Storage function path: %s", fnPath)
	}
}

// TestRemoveFunction tests the removal of a single function from the gadget.
func TestRemoveFunction(t *testing.T) {
	// Create a new Gadget instance
	gadget := &Gadget{
		Name:         "test_gadget_remove",
		IdVendor:     0x1234,
		IdProduct:    0x5678,
		SerialNumber: "123456789",
		Manufacturer: "TestManufacturer",
		Product:      "TestProduct",
	}

	// Create a configuration and add it to the gadget
	config := Config{
		Name:          "config1",
		Configuration: "Test Configuration",
		MaxPower:      "500",
		Functions:     []Function{},
	}
	gadget.Configs = append(gadget.Configs, config)

	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		return
	}

	// Add HID keyboard function
	keyboardHid := HidFunction{
		Name:         "keyboard",
		Protocol:     0x01,
		Subclass:     0x01,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}
	err := gadget.AddFunction("hid", keyboardHid)
	if err != nil {
		t.Errorf("Error adding HID keyboard function: %v", err)
	}

	// Check if HID function is correctly added
	fnPath, exists := gadget.GetFunctionPath("hid.keyboard")
	if !exists {
		t.Errorf("Expected function path for HID keyboard, but it does not exist.")
	} else {
		t.Logf("HID keyboard function path: %s", fnPath)
	}

	// Remove the HID function
	err = gadget.RemoveFunction("hid.keyboard")
	if err != nil {
		t.Errorf("Error removing HID keyboard function: %v", err)
	}

	// Check if the function was removed
	fnPath, exists = gadget.GetFunctionPath("hid.keyboard")
	if exists {
		t.Errorf("Function path for HID keyboard should have been removed, but it still exists.")
	} else {
		t.Logf("Successfully removed HID keyboard function.")
	}
}

// TestRemoveAllFunctions tests the removal of all functions from the gadget.
func TestRemoveAllFunctions(t *testing.T) {
	// Create a new Gadget instance
	gadget := &Gadget{
		Name:         "test_gadget_remove_all",
		IdVendor:     0x1234,
		IdProduct:    0x5678,
		SerialNumber: "123456789",
		Manufacturer: "TestManufacturer",
		Product:      "TestProduct",
	}

	// Create a configuration and add it to the gadget
	config := Config{
		Name:          "config1",
		Configuration: "Test Configuration",
		MaxPower:      "500",
		Functions:     []Function{},
	}
	gadget.Configs = append(gadget.Configs, config)

	if gadget.Exists() {
		gadget.Remove()
		fmt.Println("gadget already exists! Removing Now")
		return
	}

	// Add HID keyboard function
	keyboardHid := HidFunction{
		Name:         "keyboard",
		Protocol:     0x01,
		Subclass:     0x01,
		ReportLength: 8,
		Descriptor:   []byte{0x05, 0x01, 0x09, 0x06},
	}
	err := gadget.AddFunction("hid", keyboardHid)
	if err != nil {
		t.Errorf("Error adding HID keyboard function: %v", err)
	}

	// Add Mass Storage function
	massFunc := MassStorageFunction{
		Name:  "usb0",
		Stall: true,
		Luns: []MassStorageLun{
			{
				Name:          "0",
				File:          "\n", // Replace with actual device path
				Removable:     true,
				Cdrom:         true,
				Ro:            true,
				InquiryString: "Test Mass Storage",
			},
		},
	}
	err = gadget.AddFunction("mass_storage", massFunc)
	if err != nil {
		t.Errorf("Error adding Mass Storage function: %v", err)
	}

	// Check if both functions are added
	fnPath, exists := gadget.GetFunctionPath("hid.keyboard")
	if !exists {
		t.Errorf("Expected function path for HID keyboard, but it does not exist.")
	} else {
		t.Logf("HID keyboard function path: %s", fnPath)
	}

	fnPath, exists = gadget.GetFunctionPath("mass_storage.usb0")
	if !exists {
		t.Errorf("Expected function path for Mass Storage usb0, but it does not exist.")
	} else {
		t.Logf("Mass Storage function path: %s", fnPath)
	}

	// Remove all functions
	err = gadget.RemoveAllFunctions()
	if err != nil {
		t.Errorf("Error removing all functions: %v", err)
	}

	// Check if all functions were removed
	_, exists = gadget.GetFunctionPath("hid.keyboard")
	if exists {
		t.Errorf("Function path for HID keyboard should have been removed, but it still exists.")
	} else {
		t.Logf("Successfully removed HID keyboard function.")
	}

	_, exists = gadget.GetFunctionPath("mass_storage.usb0")
	if exists {
		t.Errorf("Function path for Mass Storage usb0 should have been removed, but it still exists.")
	} else {
		t.Logf("Successfully removed Mass Storage function.")
	}
}
