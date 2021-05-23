package main

import (
   "log"
   "fmt"
   "strings"
   "strconv"
   "time"
   "net"
)

type DHCPEntry struct {
   Expire time.Time
   MAC net.HardwareAddr
   IP net.IP
   Name string
   ID string
}

var dhcpLeases []DHCPEntry

func GetLeases() []DHCPEntry {
   return dhcpLeases
}

func ParseLeases(l string) []DHCPEntry {

   for i,line := range strings.Split(l, "\n") {
      c := strings.Split(line, " ")
      if len(c) < 5 {
         break
      }
      var e DHCPEntry
      for j,ent := range strings.Split(line, " ") {
         switch j {
         // Expire Time
         case 0:
           v,_ := strconv.ParseInt(ent, 10, 64)
           e.Expire = time.Unix(v,0)
           fmt.Printf("%v ", time.Until(e.Expire).Truncate(time.Second).String())
         // MAC Address
         case 1:
           v,_ := net.ParseMAC(ent)
           e.MAC = v
           fmt.Printf("%v ", v)
         // Assigned IP
         case 2:
           v := net.ParseIP(ent)
           e.IP = v
           fmt.Printf("%v ", v)
         // Name
         case 3:
           v := ent
           e.Name = v
           fmt.Printf("%v ", v)
         // DHCP Identifer If supplied
         case 4:
           v := ent
           e.ID = v
           fmt.Printf("%v\n", v)
         default:
           log.Printf("Parse Out of range Line:%d Field:%d\n", i, j)
         }
      }
      dhcpLeases = append(dhcpLeases, e)
   }
   return dhcpLeases
}
