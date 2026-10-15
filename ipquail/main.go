//  Description = ipquail website
//  Author = Theodore Baschak

package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"
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
func ipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, r.Header.Get("X-Forwarded-For")+"\n")
}

// Handles /ptr
func ptrHandler(w http.ResponseWriter, r *http.Request) {
	addr, _ := net.LookupAddr(r.Header.Get("X-Forwarded-For"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if len(addr) > 0 {
		_, _ = io.WriteString(w, addr[0]+"\n")
	} else {
		_, _ = io.WriteString(w, "none\n")
	}
}

// Handles /asn
func asnHandler(w http.ResponseWriter, r *http.Request) {
	remote_ip := r.Header.Get("X-Forwarded-For")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if IsIPv4(remote_ip) {
		ipparts := strings.Split(remote_ip, ".")
		ipparts = ipparts[:len(ipparts)-1]
		query := strings.Join(Reverse(ipparts), ".") + ".origin.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			split := strings.Split(txtrecords[0], "|")
			asn := split[0]
			_, _ = io.WriteString(w, asn+"\n")
		} else {
			_, _ = io.WriteString(w, "ERR\n")
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
			_, _ = io.WriteString(w, asn+"\n")
		} else {
			_, _ = io.WriteString(w, "ERR\n")
		}
	} else {
		_, _ = io.WriteString(w, "ERR\n")
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Accept,Content-Type,Origin")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

// Handles /api/ip
func ipApiHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"ip": r.Header.Get("X-Forwarded-For")})
}

// Handles /api/ptr
func ptrApiHandler(w http.ResponseWriter, r *http.Request) {
	addr, _ := net.LookupAddr(r.Header.Get("X-Forwarded-For"))
	if len(addr) > 0 {
		writeJSON(w, map[string]string{"ptr": addr[0]})
	} else {
		writeJSON(w, map[string]string{"ptr": "none"})
	}
}

func asnApiHandler(w http.ResponseWriter, r *http.Request) {
	remote_ip := r.Header.Get("X-Forwarded-For")
	asn := "ERR"

	if IsIPv4(remote_ip) {
		ipparts := strings.Split(remote_ip, ".")
		ipparts = ipparts[:len(ipparts)-1]
		query := strings.Join(Reverse(ipparts), ".") + ".origin.asn.cymru.com"
		txtrecords, _ := net.LookupTXT(query)

		if len(txtrecords) > 0 {
			split := strings.Split(txtrecords[0], "|")
			asn = split[0]
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
			asn = split[0]
		}
	}

	writeJSON(w, map[string]string{"asn": asn})
}

func main() {
	mux := http.NewServeMux()

	// add a route for each page you add to the site
	// make sure you create a route handler for it

	mux.HandleFunc("GET /ip", ipHandler)
	mux.HandleFunc("GET /ptr", ptrHandler)
	mux.HandleFunc("GET /asn", asnHandler)
	mux.HandleFunc("GET /api/ip", ipApiHandler)
	mux.HandleFunc("GET /api/ptr", ptrApiHandler)
	mux.HandleFunc("GET /api/asn", asnApiHandler)

	server := &http.Server{
		Addr:              ":3000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
