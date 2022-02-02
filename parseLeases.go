package main

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"log"
    "bufio"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type DHCPEntry struct {
	Expire time.Time        `json:"expire"`
	Remain time.Duration    `json:"remain"`
	MAC    net.HardwareAddr `json:"mac"`
	Info   *OuiEntry
	IP     net.IP `json:"ip"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

var dhcpLeases []DHCPEntry

func GetLeasesJson() ([]byte, error) {
	type DhcpJSON struct {
		Expire string        `json:"expire"`
		Remain string        `json:"remain"`
		Delta  time.Duration `json:"delta"`
		Mac    string        `json:"mac"`
		Info   *OuiEntry
		Ip     string `json:"ip"`
		IpInt  uint32 `json:"ipSort"`
		Name   string `json:"name"`
		Id     string `json:"id"`
        Tag    string `json:"tag"`
        Static bool `json:"static"`
	}

	dj := make([]DhcpJSON, len(dhcpLeases))

	for i, ent := range dhcpLeases {
		dj[i] = DhcpJSON{
			Expire: ent.Expire.String(),
			Remain: ent.Remain.Round(time.Second).String(),
			Delta:  ent.Remain,
			Mac:    ent.MAC.String(),
			Info:   ent.Info,
			Ip:     ent.IP.String(),
			IpInt:  ip2int(ent.IP),
			Name:   ent.Name,
			Id:     ent.ID,
            Static: false,
		}
	}

    if lookupStr("staticfile") != "" {
	    staticAddr, err := os.Open(lookupStr("staticfile"))
	    checkFatal(err, "staticfile")
        scanner := bufio.NewScanner(staticAddr)
        for scanner.Scan() {
            txt := strings.TrimSpace(scanner.Text())
            if txt == "" {
                continue
            }
            cmd := strings.Split(txt, "=")
            if strings.ToLower(cmd[0]) != "dhcp-host" {
                continue
            }
            var newStatic DhcpJSON
            val := strings.Split(cmd[1], ",")
            log.Printf("%d %q\n", len(val), val)
            i := 0
            v := val[i]
            // first entry needs to be MAC Address
            m, err := net.ParseMAC(v)
            if err != nil {
                // no support for id: ... yet
                continue
            }
            newStatic.Mac = m.String()
            newStatic.Info = lookupMac(newStatic.Mac)
            for {
                i+=1
                v = val[i]
                if net.ParseIP(v) != nil {
                    break
                }
            }
            newStatic.Ip = v
            newStatic.IpInt = ip2int(net.ParseIP(v))
            i+=1
            v = val[i]
            newStatic.Name = v
            newStatic.Id = v
            newStatic.Expire = "Never"
            newStatic.Delta,_ = time.ParseDuration("99h59m59s")
            newStatic.Remain = newStatic.Delta.String()
            newStatic.Static = true
            dj = append(dj, newStatic)
        }
    }

	return json.MarshalIndent(dj, "", "  ")
}

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

// vim: noai:ts=4:sw=4:set expandtab:
