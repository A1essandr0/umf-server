package utils

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

func CreateUUID(length int) string {
	return uuid.New().String()[24:(24+length)]
}

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		log.Printf("found ip in X-Real-IP: %s", netIP)
		return ip, nil
	}

	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			log.Printf("one of IPs in X-Forwarded-For: %s", netIP)
			return ip, nil
		}		
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		log.Printf("found ip in r.RemoteAddr: %s", netIP)
		return ip, nil
	}		

	return "", fmt.Errorf("no valid ip found")
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}