package gadget

// MassStorageFunction represents a mass storage gadget function.
type MassStorageFunction struct {
    Name  string
    Stall bool
    Luns  []MassStorageLun
}

// Ensure MassStorageFunction implements the Function interface.
var _ Function = (*MassStorageFunction)(nil)

// GadgetFunctionName returns the name of the mass storage gadget function.
func (fn *MassStorageFunction) GadgetFunctionName() string {
    return "mass_storage." + fn.Name
}

// GadgetFunctionCreate generates steps to create the mass storage gadget function.
func (fn *MassStorageFunction) GadgetFunctionCreate() (steps Steps) {
    steps.Add(Step{Action: Write, Path: "stall", Value: boolToIntStr(fn.Stall)})
    for _, lun := range fn.Luns {
        prefix := "lun." + lun.Name
        lunSteps := Steps{
            {Action: MkdirCreateOnly, Path: "", Value: ""},
        }
        lunSteps.AddSteps(lun.lunCreate())
        lunSteps.PrependPath(prefix)
        steps.AddSteps(lunSteps)
    }
    return
}

// MassStorageLun represents a logical unit number for a mass storage function.
type MassStorageLun struct {
    Name          string
    File          string
    Removable     bool
    Cdrom         bool
    Ro            bool
    InquiryString string
}

// lunCreate generates steps to create the LUN.
func (lun *MassStorageLun) lunCreate() Steps {
    return Steps{
        {Action: Write, Path: "file", Value: lun.File},
        {Action: Write, Path: "removable", Value: boolToIntStr(lun.Removable)},
        {Action: Write, Path: "cdrom", Value: boolToIntStr(lun.Cdrom)},
        {Action: Write, Path: "ro", Value: boolToIntStr(lun.Ro)},
        {Action: Write, Path: "inquiry_string", Value: lun.InquiryString},
    }
}
