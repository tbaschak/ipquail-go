//  Description = ipquail website
//  Author = Theodore Baschak
//  Version = 1.1

package main

import (
	"encoding/hex"
	"math/big"
	"net"
	"strings"

	"github.com/pilu/traffic"
)

// IsIPv4 returns true if the string passed is an IPv4 address
func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

// IsIPv6 returns true if the string passed is an IPv6 address
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// Inet6_Aton converts an IP Address (IPv4 or IPv6) net.IP object
// to a hexadecimal representation
func Inet6_Aton(ip net.IP) string {
	ipv4 := false
	if ip.To4() != nil {
		ipv4 = true
	}

	ipInt := big.NewInt(0)
	if ipv4 {
		ipInt.SetBytes(ip.To4())
		ipHex := hex.EncodeToString(ipInt.Bytes())
		return ipHex
	}

	ipInt.SetBytes(ip.To16())
	ipHex := hex.EncodeToString(ipInt.Bytes())
	return ipHex
}

// Reverse returns a reversed slice (array)
func Reverse[T any](original []T) (reversed []T) {
	reversed = make([]T, len(original))
	copy(reversed, original)

	for i := len(reversed)/2 - 1; i >= 0; i-- {
		tmp := len(reversed) - 1 - i
		reversed[i], reversed[tmp] = reversed[tmp], reversed[i]
	}

	return
}

// Handles /ip
func ipHandler(w traffic.ResponseWriter, r *traffic.Request) {
	w.Header().Add("Content-type", "text/plain")
	w.WriteText(r.Header.Get("X-Forwarded-For"))
	w.WriteText("\n")
}

// Handles /ptr
func ptrHandler(w traffic.ResponseWriter, r *traffic.Request) {
	addr, _ := net.LookupAddr(r.Header.Get("X-Forwarded-For"))
	w.Header().Add("Content-type", "text/plain")
	if len(addr) > 0 {
		w.WriteText(addr[0])
		w.WriteText("\n")
	} else {
		w.WriteText("none")
		w.WriteText("\n")
	}
}

// Handles /asn
func asnHandler(w traffic.ResponseWriter, r *traffic.Request) {
	remote_ip := r.Header.Get("X-Forwarded-For")

	w.Header().Add("Content-type", "text/plain")

	if IsIPv4(remote_ip) {
		ipparts := strings.Split(remote_ip, ".")
		ipparts = ipparts[:len(ipparts)-1]
		query := strings.Join(Reverse(ipparts), ".") + ".origin.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			split := strings.Split(txtrecords[0], "|")
			asn := split[0]
			w.WriteText(asn)
			w.WriteText("\n")
		} else {
			w.WriteText("ERR")
			w.WriteText("\n")
		}
	} else if IsIPv6(remote_ip) {
		// do IPv6 stuff
		hexstring := Inet6_Aton(net.ParseIP(remote_ip))
		ipparts := strings.Split(hexstring[:12], "")
		query := strings.Join(Reverse(ipparts), ".") + ".origin6.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			// do stuff with the first result
			split := strings.Split(txtrecords[0], "|")
			asn := split[0]
			w.WriteText(asn)
			w.WriteText("\n")
		} else {
			w.WriteText("ERR")
			w.WriteText("\n")
		}
	} else {
		w.WriteText("ERR")
		w.WriteText("\n")
	}
}

// Handles /api/ip
func ipApiHandler(w traffic.ResponseWriter, r *traffic.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,Accept,Content-Type,Origin")
	w.Header().Add("Content-type", "application/json")
	w.WriteText("{ \"ip\": \"")
	w.WriteText(r.Header.Get("X-Forwarded-For"))
	w.WriteText("\" }")
}

// Handles /api/ptr
func ptrApiHandler(w traffic.ResponseWriter, r *traffic.Request) {
	addr, _ := net.LookupAddr(r.Header.Get("X-Forwarded-For"))
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,Accept,Content-Type,Origin")
	w.Header().Add("Content-type", "application/json")
	if len(addr) > 0 {
		w.WriteText("{ \"ptr\": \"")
		w.WriteText(addr[0])
		w.WriteText("\" }")
	} else {
		w.WriteText("{ \"ptr\": \"none\" }")
	}
}

func asnApiHandler(w traffic.ResponseWriter, r *traffic.Request) {
	remote_ip := r.Header.Get("X-Forwarded-For")

	w.Header().Add("Content-type", "text/plain")

	if IsIPv4(remote_ip) {
		ipparts := strings.Split(remote_ip, ".")
		ipparts = ipparts[:len(ipparts)-1]
		query := strings.Join(Reverse(ipparts), ".") + ".origin.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			split := strings.Split(txtrecords[0], "|")
			asn := split[0]
			w.WriteText("{ \"asn\": \"" + asn + "\" }")
		} else {
			w.WriteText("{ \"asn\": \"ERR\" }")
		}
	} else if IsIPv6(remote_ip) {
		// do IPv6 stuff
		hexstring := Inet6_Aton(net.ParseIP(remote_ip))
		ipparts := strings.Split(hexstring[:12], "")
		query := strings.Join(Reverse(ipparts), ".") + ".origin6.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			// do stuff with the first result
			split := strings.Split(txtrecords[0], "|")
			asn := split[0]
			w.WriteText("{ \"asn\": \"" + asn + "\" }")
		} else {
			w.WriteText("{ \"asn\": \"ERR\" }")
		}
	} else {
		w.WriteText("{ \"asn\": \"ERR\" }")
	}
}

func main() {
	router := traffic.New()

	// add a route for each page you add to the site
	// make sure you create a route handler for it

	//  router.Get("/", indexHandler)
	router.Get("/ip", ipHandler)
	router.Get("/ptr", ptrHandler)
	router.Get("/asn", asnHandler)
	router.Get("/api/ip", ipApiHandler)
	router.Get("/api/ptr", ptrApiHandler)
	router.Get("/api/asn", asnApiHandler)
	router.Run()
}
