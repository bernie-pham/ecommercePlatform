package ultils

func GetOptionalInt(value int) (int, bool) {
	if value > 0 {
		return value, true
	}
	return 0, false
}
