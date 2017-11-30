package main

import (
	"github.com/xlab/portmidi"
)


func kashmirMangle() filterFunc {
	return func(s, d *Device, event *portmidi.Event) bool {
		if program != 0x1900 {
			return false
		}

		if ok, v := nordSlotChangeMessage(s, event); ok {
			switch v {
			case 0x00:
				setLeadHold(d, true, event.Timestamp)
			default:
				setLeadHold(d, false, event.Timestamp)
			}
		}

		switch slot {
		case 0x00:
			return false
		case 0x2b:
			return s.matches("Electro") && d.matches("Lead") // slot B, electro plays lead
		case 0x55:
			return s.matches("Lead") && d.matches("Electro") // slot C, lead plays electro
		case 0x7f:
			return true
		default:
			return false
		}

		return false
	}
}
