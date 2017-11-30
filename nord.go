package main

import (
	"github.com/xlab/portmidi"
)

var (
	slot byte
)

const (
	STATUS_CC   = 0xb0
	STATUS_PROG = 0xc0

	CC_BANK_MSB    = 0x00
	CC_BANK_LSB    = 0x20
	CC_SLOT_CHANGE = 0x31
	CC_HOLD        = 0x3a

	HOLD_OFF = 0x00
	HOLD_ON  = 0x7f

	SLOT_A = 0x00
	SLOT_B = 0x2b
	SLOT_C = 0x55
	SLOT_D = 0x7f
)

func nordSlotChangeMessage(s *Device, event *portmidi.Event) (bool, byte) {
	if s.matches("Electro") && event.Message.Status() == STATUS_CC && event.Message.Data1() == CC_SLOT_CHANGE {
		return true, event.Message.Data2()
	}
	return false, 0xFF
}

func nordElectroSlotTracking() filterFunc {
	return func(s, d *Device, event *portmidi.Event) bool {
		if ok, v := nordSlotChangeMessage(s, event); ok {
			slot = v
		}

		return false
	}
}

func setLeadHold(d *Device, state bool, timestamp int32) {
	var val byte = HOLD_OFF
	if state {
		val = HOLD_ON
	}
	d.outputStream.Sink() <- portmidi.Event{timestamp, portmidi.NewMessage(STATUS_CC, CC_HOLD, val), nil}
}

/*
	expressionRewrite := func(s, d *Device, event *portmidi.Event) bool {
		if strings.Contains(s.Name, "Electro") && strings.Contains(d.Name, "Lead") {
			if event.Message.Status() == 0xb0 && event.Message.Data1() == 11 { // CC #11
				// rewrite to CC #7
				event.Message = portmidi.NewMessage(0xb0, 7, event.Message.Data2())
				return true
			}
		}
		return false
	}
*/
