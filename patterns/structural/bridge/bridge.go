package main

import "fmt"

// Когда использовать: Когда нужно разделить абстракцию и реализацию так, чтобы они могли изменяться независимо.

// Интерфейс устройства
type Device interface {
	isEnabled() bool
	enable()
	disable()
	getVolume() int
	setVolume(int)
}

// Конкретные устройства
type Radio struct {
	volume  int
	enabled bool
}

func (r *Radio) isEnabled() bool {
	return r.enabled
}

func (r *Radio) enable() {
	r.enabled = true
}

func (r *Radio) disable() {
	r.enabled = false
}

func (r *Radio) getVolume() int {
	return r.volume
}

func (r *Radio) setVolume(volume int) {
	r.volume = volume
}

// Абстракция
type Remote struct {
	device Device
}

func (r *Remote) TogglePower() {
	if r.device.isEnabled() {
		r.device.disable()
	} else {
		r.device.enable()
	}
}

func (r *Remote) VolumeDown() {
	r.device.setVolume(r.device.getVolume() - 10)
}

func (r *Remote) VolumeUp() {
	r.device.setVolume(r.device.getVolume() + 10)
}

func main() {
	radio := &Radio{}
	remote := &Remote{device: radio}
	remote.TogglePower()
	fmt.Println("Радио включено:", radio.isEnabled())
}
