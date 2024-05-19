package mygo

import "testing"

func TestLocalIP(t *testing.T) {
	ip, err := Net.LocalIP()

	if err != nil {
		t.Error(err)
	}

	t.Log(ip)
}

func TestLocalMAC(t *testing.T) {
	mac, err := Net.LocalMac()

	if err != nil {
		t.Error(err)
	}

	t.Log(mac)
}
