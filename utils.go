package gadget

// Convert boolean to "1" or "0" string representation.
func boolToIntStr(b bool) string {
    if b {
        return "1"
    }
    return "0"
}

// GetUdcs returns a list of UDCs (USB Device Controllers).
func GetUdcs() []string {
    var udcs []string

    files, err := os.ReadDir(UdcPathGlob)
    if err != nil {
        return nil
    }

    for _, file := range files {
        udcs = append(udcs, file.Name())
    }

    return udcs
}