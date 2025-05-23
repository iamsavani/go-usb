package gadget

// MassStorageFunction represents a mass storage gadget function.
type MassStorageFunction struct {
	Name    string
	Stall   bool
	Luns    []MassStorageLun
	Enabled bool
}

// Ensure MassStorageFunction implements the Function interface.
var _ Function = (*MassStorageFunction)(nil)

// GadgetFunctionName returns the name of the mass storage gadget function.
func (fn *MassStorageFunction) GadgetFunctionName() string {
	return fn.Name
}

// GadgetFunctionType returns the function type name (without the instance).
func (fn *MassStorageFunction) GadgetFunctionType() string {
	return "mass_storage"
}

// GadgetFunctionCreate generates steps to create the mass storage gadget function.
func (fn *MassStorageFunction) GadgetFunctionCreate() (steps Steps) {
	steps.Append(Step{Write, "stall", boolToIntStr(fn.Stall)})
	for _, lun := range fn.Luns {
		var (
			prefix = "lun." + lun.Name
			lsteps = Steps{
				Step{MkdirCreateOnly, "", ""},
			}
		)

		lsteps.Extend(lun.lunCreate())
		lsteps.PrependPath(prefix)

		steps.Extend(lsteps)
	}
	return
}

// MassStorageLun represents a logical unit number for a mass storage function.
type MassStorageLun struct {
	Name string
	Attr GadgetAttributes
}

// lunCreate generates steps to create the LUN.
func (lun *MassStorageLun) lunCreate() Steps {
	steps := Steps{}
	for key, value := range lun.Attr {
		steps = append(steps, Step{Write, key, value})
	}
	return steps
}

func (fn *MassStorageFunction) isEnabled() bool {
	return fn.Enabled
}
