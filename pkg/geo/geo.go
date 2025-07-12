package geo

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	GEOIP_PROTO = getEnv("GEOIP_PROTO", "http")
	GEOIP_HOST  = getEnv("GEOIP_HOST", "localhost:8080")
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type GeoInfo struct {
	IP         string `json:"ip"`
	Country    string `json:"country"`
	CountryISO string `json:"country_iso"`
	RegionName string `json:"region_name"`
	RegionCode string `json:"region_code"`
	City       string `json:"city"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

func IPFromRequest(headers []string, r *http.Request) (net.IP, error) {
	remoteIP := ""
	for _, h := range headers {
		remoteIP = r.Header.Get(h)
		if http.CanonicalHeaderKey(h) == "X-Forwarded-For" {
			remoteIP = ipFromForwardedForHeader(remoteIP)
		}
		if remoteIP != "" {
			break
		}
	}

	if remoteIP == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		remoteIP = host
	}

	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP %s", remoteIP)
	}
	return ip, nil
}

func ipFromForwardedForHeader(v string) string {
	sep := strings.Index(v, ",")
	if sep == -1 {
		return v
	}
	return v[:sep]
}

func GetGeoInfo(ip string) (*GeoInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/json?ip=%s", GEOIP_PROTO, GEOIP_HOST, ip), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info GeoInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	return &info, err
}
