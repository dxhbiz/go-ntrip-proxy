package proxy

import (
	"sync"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/config"
)

type Action int

var (
	casters map[string]*Caster
)

const (
	TickAction Action = iota
	SendAction
)

type SubMsg struct {
	action Action
	data   []byte
}

type SubClient interface {
	IsMain() bool
	OnData(msg SubMsg)
}

type Caster struct {
	cfg    config.CasterConfig
	subs   sync.Map
	puller *Puller
}

// init casters
func initCasters() {
	cfg := config.GetConfig()

	casters = make(map[string]*Caster)

	cfgCasters := cfg.Casters

	for _, caster := range cfgCasters {
		casters[caster.Name] = newCaster(caster)
	}
}

func newCaster(casterCfg config.CasterConfig) *Caster {
	caster := &Caster{
		cfg:    casterCfg,
		puller: newPuller(casterCfg),
	}
	return caster
}

func (c *Caster) addMain() {
	c.subs.Range(func(key, value any) bool {
		k := key.(string)
		nc := value.(*ntripClient)
		if nc.IsMain() {
			key = k

			msg := SubMsg{
				action: TickAction,
			}
			nc.OnData(msg)

			c.subs.Delete(key)
			return false
		}
		return true
	})

	c.puller.changeSignal(ConnectSignal)
}

func (c *Caster) addSub(id string, sub SubClient) bool {
	c.subs.Store(id, sub)
	return true
}

func (c *Caster) delSub(id string) {
	c.subs.Delete(id)
}

func (c *Caster) broadcast(data []byte) {
	c.subs.Range(func(key, value any) bool {
		nc := value.(*ntripClient)
		msg := SubMsg{
			action: SendAction,
			data:   data,
		}
		nc.OnData(msg)

		return true
	})
}

func (c *Caster) sendNmea(data []byte) {
	c.puller.send(data)
}

func subscribe(mountpoint, id string, isMain bool, sub *ntripClient) bool {
	if caster, ok := casters[mountpoint]; ok {
		if isMain {
			caster.addMain()
		}
		return caster.addSub(id, sub)
	} else {
		return false
	}
}

func unSubscribe(mountpoint, id string, isMain bool) {
	if caster, ok := casters[mountpoint]; ok {
		if isMain {
			caster.puller.changeSignal(CloseSignal)
		}

		caster.delSub(id)
	}
}

func publish(mountpoint string, data []byte) {
	if caster, ok := casters[mountpoint]; ok {
		caster.broadcast(data)
	}
}

func sendNmea(mountpoint string, nmea []byte) {
	if caster, ok := casters[mountpoint]; ok {
		caster.sendNmea(nmea)
	}
}
