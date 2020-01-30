package gintunnel

import "sync"

func NewTunnel() {
	var wg sync.WaitGroup
	wg.Add(2)
	rg := Register{fm: NewForwardMap()}
	fw := Forwarder{register: rg}
	go func() {
		rg.Start()
		wg.Done()
	}()
	go func() {
		fw.Start()
		wg.Done()
	}()
	wg.Wait()
}
