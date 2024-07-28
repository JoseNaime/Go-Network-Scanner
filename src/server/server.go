package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/gopacket/pcap"
	"go-network-scanner-nmap/src/network"
	"go-network-scanner-nmap/src/nmaputil"
	"log"
	"net"
	"time"
)

type ScanResult struct {
	MacAddress   string                    `json:"mac_address"`
	IPAddress    string                    `json:"ip_address"`
	Mask         string                    `json:"mask_address"`
	NmapScanData []nmaputil.NmapScanResult `json:"nmap_scan_data"`
	Timestamp    time.Time                 `json:"timestamp"`
}

type ServerScanResponse struct {
	Status string      `json:"status"`
	Time   interface{} `json:"time"`
	Data   interface{} `json:"data"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan map[string]interface{})
var scanResult = ScanResult{}

func StartServer(nmapPath string, ifaceName *string) {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			delete(clients, c)
			c.Close()
		}()
		clients[c] = true
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
			log.Printf("Received message: %s", msg)
		}
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	app.Get("/stored", func(c *fiber.Ctx) error {
		log.Println("Stored info requested")
		res := ServerScanResponse{}

		if scanResult.IPAddress == "" {
			res = ServerScanResponse{
				Status: "Scanning...",
				Time:   time.Now(),
				Data:   nil,
			}
		} else {
			res = ServerScanResponse{
				Status: "Stored",
				Time:   scanResult.Timestamp,
				Data:   scanResult,
			}
		}

		return c.JSON(res)
	})

	go runNetworkScan(nmapPath, ifaceName)
	go handleMessages()

	err := app.Listen("0.0.0.0:5234")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func runNetworkScan(nmapPath string, ifaceName *string) {
	for {
		fmt.Println("Scanning network...    " + time.Now().String())
		handle, err := pcap.OpenLive(*ifaceName, 65536, true, pcap.BlockForever)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		iface, err := net.InterfaceByName(*ifaceName)
		if err != nil {
			log.Fatal(err)
		}

		ip, ipNet, err := network.GetInterfaceIPv4Addr(*ifaceName)
		if err != nil {
			log.Fatal(err)
		}

		activeIps := network.GetActiveIPsInSubnet(ipNet)
		fmt.Printf("Amount of active IPs: %d\n", len(activeIps))

		scanData := nmaputil.ScanIpsWithNmap(activeIps, nmapPath)

		scanResult = ScanResult{
			MacAddress:   iface.HardwareAddr.String(),
			IPAddress:    ip.String(),
			Mask:         ipNet.Mask.String(),
			NmapScanData: scanData,
			Timestamp:    time.Now(),
		}

		broadcast <- map[string]interface{}{
			"status": "Realtime Data",
			"time":   time.Now(),
			"data":   scanResult,
		}

		// Delay between scans
		time.Sleep(1 * time.Minute)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing JSON to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
