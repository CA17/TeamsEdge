package sysid

import (
	"crypto/md5"
	"encoding/hex"
	"net"

	"github.com/denisbrodbeck/machineid"
)

func GetSystemSid() string {
	sid, err := machineid.ProtectedID("TeamsEdge")
	if err != nil {
		return GetFirstMacAddrSid()
	}else{
		return EncodeHash(sid)
	}
}

func GetFirstMacAddrSid() string {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.HardwareAddr.String() != "" {
			return EncodeHash(inter.HardwareAddr.String())
		}
	}
	return EncodeHash("TeamsEdge")
}

func EncodeHash(src string) string {
	hash := md5.Sum([]byte(src))
	return hex.EncodeToString(hash[:])
}

