package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/pion/stun"
	"github.com/pion/turn/v2"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func toInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("Can't convert %s to int", value)
	}
	return result
}

// stunLogger wraps a PacketConn and prints incoming/outgoing STUN packets
// This pattern could be used to capture/inspect/modify data as well
type stunLogger struct {
	net.PacketConn
}

func (s *stunLogger) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if n, err = s.PacketConn.WriteTo(p, addr); err == nil && stun.IsMessage(p) {
		msg := &stun.Message{Raw: p}
		if err = msg.Decode(); err != nil {
			return
		}

		fmt.Printf("Outbound STUN: %s \n", msg.String())
	}

	return
}

func (s *stunLogger) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	if n, addr, err = s.PacketConn.ReadFrom(p); err == nil && stun.IsMessage(p) {
		msg := &stun.Message{Raw: p}
		if err = msg.Decode(); err != nil {
			return
		}

		fmt.Printf("Inbound STUN: %s \n", msg.String())
	}

	return
}

func main() {
	publicIP := flag.String("ip", getEnv("GRUPPETTO_IP", "127.0.0.1"), "IP Address that TURN can be contacted by.")
	port := flag.Int("port", toInt(getEnv("GRUPPETTO_PORT", "3478")), "Listening port.")
	user := flag.String("user", getEnv("GRUPPETTO_USER", "user"), "Username.")
	password := flag.String("password", getEnv("GRUPPETTO_PASSWORD", "password"), "Password.")
	realm := flag.String("realm", getEnv("GRUPPETTO_REALM", "gruppetto"), "Realm.")
	flag.Parse()

	if len(*publicIP) == 0 {
		log.Panicf("'ip' is required")
	} else if len(*user) == 0 {
		log.Panicf("'user' is required")
	} else if len(*password) == 0 {
		log.Panicf("'password' is required")
	}

	// Create a UDP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	// Create a TCP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any TCP listeners, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	tcpListener, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	// Cache user for easy lookup later
	// If passwords are stored they should be saved to your DB hashed using turn.GenerateAuthKey
	usersMap := map[string][]byte{
		*user: turn.GenerateAuthKey(*user, *realm, *password),
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: *realm, // Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		}, // PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: &stunLogger{udpListener},
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(*publicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",              // But actually be listening on every interface
				},
			},
		}, // ListenerConfig is a list of Listeners and the configuration around them
		ListenerConfigs: []turn.ListenerConfig{
			{
				Listener: tcpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(*publicIP),
					Address:      "0.0.0.0",
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}

	log.Print("Server started.")

	// Block until user sends SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Print("Closing server...")

	if err = s.Close(); err != nil {
		log.Panic(err)
	}

	log.Print("Server closed.")
}
