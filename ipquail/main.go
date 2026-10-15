//  Description = ipquail website
//  Author = Theodore Baschak

package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

const realIPHeader = "X-Real-IP"

func requestIP(w http.ResponseWriter, r *http.Request) (netip.Addr, bool) {
	peer, err := netip.ParseAddrPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid remote address", http.StatusBadRequest)
		return netip.Addr{}, false
	}
	if !peer.Addr().IsLoopback() {
		http.Error(w, "forbidden", http.StatusForbidden)
		return netip.Addr{}, false
	}

	ip, err := netip.ParseAddr(r.Header.Get(realIPHeader))
	if err != nil {
		http.Error(w, "invalid client address", http.StatusBadRequest)
		return netip.Addr{}, false
	}

	return ip.Unmap(), true
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

func asnQuery(ip netip.Addr) string {
	if !ip.IsValid() {
		return ""
	}
	if ip.Is4() {
		octets := ip.As4()
		parts := []string{
			strconv.Itoa(int(octets[2])),
			strconv.Itoa(int(octets[1])),
			strconv.Itoa(int(octets[0])),
		}
		return strings.Join(parts, ".") + ".origin.asn.cymru.com"
	}

	octets := ip.As16()
	hexString := hex.EncodeToString(octets[:])
	nibbles := strings.Split(hexString[:16], "") // Aggregate IPv6 lookups at /64.
	return strings.Join(Reverse(nibbles), ".") + ".origin6.asn.cymru.com"
}

// Handles /ip
func ipHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, remoteIP.String()+"\n")
}

// Handles /ptr
func ptrHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	addr, _ := net.LookupAddr(remoteIP.String())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if len(addr) > 0 {
		_, _ = io.WriteString(w, addr[0]+"\n")
	} else {
		_, _ = io.WriteString(w, "none\n")
	}
}

// Handles /asn
func asnHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	txtrecords, _ := net.LookupTXT(asnQuery(remoteIP))
	if len(txtrecords) > 0 {
		split := strings.Split(txtrecords[0], "|")
		asn := split[0]
		_, _ = io.WriteString(w, asn+"\n")
	} else {
		_, _ = io.WriteString(w, "ERR\n")
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

// Handles /api/ip
func ipApiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	writeJSON(w, map[string]string{"ip": remoteIP.String()})
}

// Handles /api/ptr
func ptrApiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	addr, _ := net.LookupAddr(remoteIP.String())
	if len(addr) > 0 {
		writeJSON(w, map[string]string{"ptr": addr[0]})
	} else {
		writeJSON(w, map[string]string{"ptr": "none"})
	}
}

func asnApiHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP, ok := requestIP(w, r)
	if !ok {
		return
	}
	asn := "ERR"
	txtrecords, _ := net.LookupTXT(asnQuery(remoteIP))
	if len(txtrecords) > 0 {
		split := strings.Split(txtrecords[0], "|")
		asn = split[0]
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
		Addr:              "127.0.0.1:3000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
