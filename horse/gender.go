package horse

type Gender int

const (
	GenderInvalid Gender = iota
	GenderStallion
	GenderGelding
	GenderMare
)

// A horse is a filly or colt when they are less than 4 years old
const maxYouthAge = 3

func (h *Horse) Gender() string {
	isYouth := h.Age() < maxYouthAge
	switch h.gender {
	case GenderStallion:
		if isYouth {
			return "Colt"
		}
		return "Stallion"
	case GenderGelding:
		if isYouth {
			return "Colt"
		}
		return "Gelding"
	case GenderMare:
		if isYouth {
			return "Filly"
		}
		return "Mare"
	}
	return ""
}
