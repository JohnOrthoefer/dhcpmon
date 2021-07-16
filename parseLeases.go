package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type DHCPEntry struct {
	Expire  time.Time
	Remain  time.Duration
	MAC     net.HardwareAddr
	Info    *OuiEntry
	IP      net.IP
	Name    string
	ID      string
}

var dhcpLeases []DHCPEntry

func GetLeases() []DHCPEntry {
	for i, ent := range dhcpLeases {
		dhcpLeases[i].Remain = time.Until(ent.Expire).Truncate(time.Second)
	}
	return dhcpLeases
}

func readLeases(f string) {
	leases, err := os.ReadFile(f)
	checkFatal(err, f)

	ParseLeases(string(leases))
}

func watchLeases(f string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modiied file", event.Name)
					readLeases(f)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(f)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func ParseLeases(l string) []DHCPEntry {
	dhcpLeases = dhcpLeases[:0]

	for i, line := range strings.Split(l, "\n") {
		c := strings.Split(line, " ")
		if len(c) < 5 {
			break
		}
		var e DHCPEntry
		for j, ent := range strings.Split(line, " ") {
			switch j {
			// Expire Time
			case 0:
				v, _ := strconv.ParseInt(ent, 10, 64)
				e.Expire = time.Unix(v, 0)
				e.Remain = time.Until(e.Expire)
				//            fmt.Printf("%v ", time.Until(e.Expire).Truncate(time.Second).String())
			// MAC Address
			case 1:
				v, _ := net.ParseMAC(ent)
				e.Info = lookupMac(ent)
				log.Printf("Mac Info(%v): %v", v, getCompany(e.Info))
				e.MAC = v
				//            fmt.Printf("%v ", v)
			// Assigned IP
			case 2:
				v := net.ParseIP(ent)
				e.IP = v
				//            fmt.Printf("%v ", v)
			// Name
			case 3:
				v := ent
				e.Name = v
				//            fmt.Printf("%v ", v)
			// DHCP Identifer If supplied
			case 4:
				v := ent
				e.ID = v
				//            fmt.Printf("%v\n", v)
			default:
				log.Printf("Parse Out of range Line:%d Field:%d\n", i, j)
			}
		}
		dhcpLeases = append(dhcpLeases, e)
	}
	return dhcpLeases
}
