package gadget

import "strconv"

// HidFunction represents a HID gadget function.
type HidFunction struct {
	Name         string
	Protocol     uint8
	Subclass     uint8
	ReportLength uint16
	Descriptor   []byte
}

// Ensure HidFunction implements the Function interface.
var _ Function = (*HidFunction)(nil)

// GadgetFunctionName returns the name of the HID gadget function.
func (fn *HidFunction) GadgetFunctionName() string {
	return "hid." + fn.Name
}

// GadgetFunctionCreate generates steps to create the HID gadget function.
func (fn *HidFunction) GadgetFunctionCreate() Steps {
	return Steps{
		Step{Write, "protocol", strconv.Itoa(int(fn.Protocol))},
		Step{Write, "subclass", strconv.Itoa(int(fn.Subclass))},
		Step{Write, "report_length", strconv.Itoa(int(fn.ReportLength))},
		Step{WriteBinary, "report_desc", string(fn.Descriptor)},
	}
}
