package main

import (
	"github.com/xlab/portmidi"
)

const (
	PROG_KASHMIR = 0x1900
)


func kashmirMangle() filterFunc {
	return func(s, d *Device, event *portmidi.Event) bool {
		if program != PROG_KASHMIR {
			return false
		}

		if ok, v := nordSlotChangeMessage(s, event); ok {
			switch v {
			case SLOT_A:
				setLeadHold(d, true, event.Timestamp)
			default:
				setLeadHold(d, false, event.Timestamp)
			}
		}

		switch slot {
		case SLOT_A:
			return false
		case SLOT_B:
			return s.matches("Electro") && d.matches("Lead") // slot B, electro plays lead
		case SLOT_C:
			return s.matches("Lead") && d.matches("Electro") // slot C, lead plays electro
		case SLOT_D:
			return true
		default:
			return false
		}

		return false
	}
}
