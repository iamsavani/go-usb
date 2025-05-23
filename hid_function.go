package gadget

// HidFunction represents a HID gadget function.
type HidFunction struct {
	Name       string
	Attrs      GadgetAttributes
	ReportDesc []byte
	Enabled    bool
}

// Ensure HidFunction implements the Function interface.
var _ Function = (*HidFunction)(nil)

// GadgetFunctionName returns the name of the HID gadget function.
func (fn *HidFunction) GadgetFunctionName() string {
	return fn.Name
}

// GadgetFunctionType returns the function type name (without the instance).
func (fn *HidFunction) GadgetFunctionType() string {
	return "hid"
}

// GadgetFunctionCreate generates steps to create the HID gadget function.
func (fn *HidFunction) GadgetFunctionCreate() Steps {
	steps := Steps{}
	for key, value := range fn.Attrs {
		steps = append(steps, Step{Write, key, value})
	}
	steps = append(steps, Step{WriteBinary, "report_desc", string(fn.ReportDesc)})
	return steps
}

func (fn *HidFunction) isEnabled() bool {
	return fn.Enabled
}
