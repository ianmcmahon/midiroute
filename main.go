package main

import (
	"fmt"
	"strings"

	"github.com/xlab/closer"
	"github.com/xlab/portmidi"
)

type Device struct {
	Name string
	inputDeviceID portmidi.DeviceID
	outputDeviceID portmidi.DeviceID

	inputStream *portmidi.Stream
	outputStream *portmidi.Stream
}

func (d *Device) matches(s string) bool {
	return strings.Contains(d.Name, s)
}

type filterFunc func(s, d *Device, event *portmidi.Event) bool

var (
	program uint32 = 0x00000000
)

func programChangeFilter() filterFunc {
	return func(s, d *Device, event *portmidi.Event) bool {
		if event.Message.Status() == 0xc0 {
			program &= 0xFFFF0000
			program |= uint32(event.Message.Data1()) << 8
			program |= uint32(event.Message.Data2()) << 16
			return true
		}
		if event.Message.Status() == 0xb0 {
			if event.Message.Data1() == 0x00 {
				program &= 0x00FFFFFF
				program |= uint32(event.Message.Data2()) << 24
				return true
			}
			if event.Message.Data1() == 0x20 {
				program &= 0xFF00FFFF
				program |= uint32(event.Message.Data2()) << 16
				return true
			}
		}
		return false
	}
}

func devicesMatching(match ...string) map[string]*Device {
	numDevices := portmidi.CountDevices()
	fmt.Printf("total devices: %d\n", numDevices)

	devices := map[string]*Device{}

	for i := 0; i < numDevices; i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			fmt.Printf("\t[%s]\n", info.Name)
		}
		device, ok := devices[info.Name]
		if !ok {
			device = &Device{Name: info.Name, inputDeviceID: -1, outputDeviceID: -1}
		}
		if info.IsInputAvailable {
			device.inputDeviceID = portmidi.DeviceID(i)
		}
		if info.IsOutputAvailable {
			device.outputDeviceID = portmidi.DeviceID(i)
		}

		if len(match) == 0 {
			devices[device.Name] = device
		} else {
			for _, m := range match {
				if device.matches(m) {
					devices[device.Name] = device
					continue
				}
			}
		}
	}

	return devices
}

func main() {
	defer closer.Close()

	portmidi.Initialize()
	closer.Bind(func() {
		portmidi.Terminate()
	})

	devices := devicesMatching("Nord")

	for n, dev := range devices {
		fmt.Printf("Opening device: %s\n", n)
		var err error
		dev.inputStream, err = portmidi.NewInputStream(dev.inputDeviceID, 1024, 0)
		if err != nil {
			closer.Fatalln("[ERR] cannot init input stream: ", err)
		}
		closer.Bind(func() {
			dev.inputStream.Close()
		})
		dev.outputStream, err = portmidi.NewOutputStream(dev.outputDeviceID, 1024, 0, 0)
		if err != nil {
			closer.Fatalln("[ERR] cannot init output stream: ", err)
		}
		closer.Bind(func() {
			dev.outputStream.Close()
		})
	}

	// set up filters and routing rules
	// note these are positive filters, meaning returning true means forward the message

	/*
	*/

	filters := []filterFunc{
		programChangeFilter(),
		nordElectroSlotTracking(),
		kashmirMangle(),
	}

	// now that all streams are open, we need to do an NxN route matrix

	for _, s := range devices {
		go func(src *Device) {
			for ev := range src.inputStream.Source() {
				for _, dest := range devices {
					if src == dest {
						continue
					}
					for _, filter := range filters {
						if filter(src, dest, &ev) {
							//fmt.Printf("%s -> %s: %x %x %x\n", src.Name, dest.Name, ev.Message.Status(), ev.Message.Data1(), ev.Message.Data2())
							dest.outputStream.Sink() <- ev
						}
					}
				}
			}
		}(s)
	}

	closer.Hold()
}
