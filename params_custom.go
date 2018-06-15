package llrp

func CustomParameter(VendorID, SubType int, params ...interface{}) []interface{} {
	inf := []interface{}{
		uint32(VendorID),
		uint32(SubType),
	}
	for _, k := range params {
		inf = append(inf, k)
	}
	return commonSpec(
		P_Custom,
		inf,
	)
}
