// Copyright (C) 2018 LEAP
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package systray

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"0xacab.org/leap/bitmask-vpn/pkg/bitmask"
	"0xacab.org/leap/bitmask-vpn/pkg/config"
)

type bmTray struct {
	bm            bitmask.Bitmask
	conf          *Config
	waitCh        chan bool
	activeGateway *gatewayTray
	autostart     bitmask.Autostart
}

type gatewayTray struct {
	name string
}

func (bt *bmTray) start() {
	// XXX this removes the snap error message, but produces an invisible icon.
	// https://0xacab.org/leap/riseup_vpn/issues/44
	// os.Setenv("TMPDIR", "/var/tmp")
}

func (bt *bmTray) onExit() {
	log.Println("Closing systray")
}

func (bt *bmTray) onReady() {
	bt.waitCh <- true
}

func (bt *bmTray) setUpSystray() {

	if bt.conf.SelectGateway {
		bt.addGateways()
	}

	showDonate, err := strconv.ParseBool(config.AskForDonations)
	if err != nil {
		log.Printf("Error parsing AskForDonations: %v", err)
		showDonate = true
	}
	if !showDonate {
	}

}

func (bt *bmTray) loop(bm bitmask.Bitmask, as bitmask.Autostart) {
	<-bt.waitCh
	bt.waitCh = nil

	bt.bm = bm
	bt.autostart = as
	bt.setUpSystray()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	if status, err := bt.bm.GetStatus(); err != nil {
		log.Printf("Error getting status: %v", err)
	} else {
		bt.changeStatus(status)
	}

}

func (bt *bmTray) addGateways() {
	_, err := bt.bm.ListGateways(config.Provider)
	//gatewayList, err := bt.bm.ListGateways(config.Provider)
	if err != nil {
		log.Printf("Gateway initialization error: %v", err)
		return
	}

	/*
		for i, city := range gatewayList {
			menuItem := systray.AddMenuItem(city, bt.conf.Printer.Sprintf("Use %s %v gateway", config.ApplicationName, city))
			gateway := gatewayTray{menuItem, city}

			if i == 0 {
				menuItem.Check()
				menuItem.SetTitle("*" + city)
				bt.activeGateway = &gateway
			} else {
				menuItem.Uncheck()
			}

			go func(gateway gatewayTray) {
				for {
					<-menuItem.ClickedCh
					gateway.menuItem.SetTitle("*" + gateway.name)
					gateway.menuItem.Check()

					bt.activeGateway.menuItem.Uncheck()
					bt.activeGateway.menuItem.SetTitle(bt.activeGateway.name)
					bt.activeGateway = &gateway

					bt.bm.UseGateway(gateway.name)
					log.Printf("Manual connection to %s gateway\n", gateway.name)
					bt.bm.StartVPN(config.Provider)
				}
			}(gateway)
		}
	*/
}

func (bt *bmTray) changeStatus(status string) {
	printer := bt.conf.Printer
	statusStr := ""

	switch status {
	case "on":
		statusStr = printer.Sprintf("%s on", config.ApplicationName)

	case "off":
		statusStr = printer.Sprintf("%s off", config.ApplicationName)

	case "starting":
		bt.waitCh = make(chan bool)
		go bt.waitIcon()
		statusStr = printer.Sprintf("Connecting to %s", config.ApplicationName)

	case "stopping":
		bt.waitCh = make(chan bool)
		go bt.waitIcon()
		statusStr = printer.Sprintf("Stopping %s", config.ApplicationName)

	case "failed":
		statusStr = printer.Sprintf("%s blocking internet", config.ApplicationName)
	}
	fmt.Println(statusStr)
}

func (bt *bmTray) waitIcon() {
	//icons := [][]byte{icon.Wait0, icon.Wait1, icon.Wait2, icon.Wait3}
	//_ := [][]byte{icon.Wait0, icon.Wait1, icon.Wait2, icon.Wait3}
	for i := 0; true; i = (i + 1) % 4 {
		select {
		case <-bt.waitCh:
			return
		case <-time.After(time.Millisecond * 500):
			continue
		}
	}
}
