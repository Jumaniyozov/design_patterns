package bridge

import "fmt"

type Device interface {
	IsEnabled() bool
	Enable()
	Disable()
	GetVolume() int
	SetVolume(percent int)
}

type TV struct {
	on     bool
	volume int
}

func (t *TV) IsEnabled() bool { return t.on }
func (t *TV) Enable()         { t.on = true; fmt.Println("TV: turned on") }
func (t *TV) Disable()        { t.on = false; fmt.Println("TV: turned off") }
func (t *TV) GetVolume() int  { return t.volume }
func (t *TV) SetVolume(percent int) {
	t.volume = percent
	fmt.Printf("TV: volume set to %d%%\n", percent)
}

type Radio struct {
	on     bool
	volume int
}

func (r *Radio) IsEnabled() bool { return r.on }
func (r *Radio) Enable()         { r.on = true; fmt.Println("Radio: turned on") }
func (r *Radio) Disable()        { r.on = false; fmt.Println("Radio: turned off") }
func (r *Radio) GetVolume() int  { return r.volume }
func (r *Radio) SetVolume(percent int) {
	r.volume = percent
	fmt.Printf("Radio: volume set to %d%%\n", percent)
}

type RemoteControl struct {
	device Device
}

func NewRemoteControl(device Device) *RemoteControl {
	return &RemoteControl{device: device}
}

func (r *RemoteControl) TogglePower() {
	if r.device.IsEnabled() {
		r.device.Disable()
	} else {
		r.device.Enable()
	}
}

func (r *RemoteControl) VolumeUp() {
	r.device.SetVolume(r.device.GetVolume() + 10)
}

func (r *RemoteControl) VolumeDown() {
	r.device.SetVolume(r.device.GetVolume() - 10)
}
