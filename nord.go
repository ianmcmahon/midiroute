package main

import (
	"github.com/xlab/portmidi"
)

var (
	slot byte
)

func nordSlotChangeMessage(s *Device, event *portmidi.Event) (bool, byte) {
	if s.matches("Electro") && event.Message.Status() == 0xb0 && event.Message.Data1() == 0x31 {
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
	var val byte = 0x00
	if state {
		val = 0x7f
	}
	d.outputStream.Sink() <- portmidi.Event{timestamp, portmidi.NewMessage(0xb0, 0x3a, val), nil}
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
