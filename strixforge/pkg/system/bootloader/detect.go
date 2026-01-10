package bootloader

// Detect returns a list of all detected active bootloaders
func Detect() []Bootloader {
	candidates := []Bootloader{
		NewGrub(),
		NewSystemdBoot(),
		NewLimine(),
		NewRefind(),
	}

	active := []Bootloader{}
	for _, b := range candidates {
		if b.IsInstalled() {
			active = append(active, b)
		}
	}
	return active
}
