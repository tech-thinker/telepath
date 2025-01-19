package handler

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/utils"
)

func (h *handler) addTunnel(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToTunnel(packet.Data)
	if err != nil {
		log.Println("Uanble to decode tunnel")
		return packet, err
	}
	err = h.configRepo.AddTunnel(doc.Name, doc.LocalPort, doc.RemoteHost, doc.RemotePort, doc.HostChain)
	if err != nil {
		log.Println("failed to add tunnel.", err)
	}
	result.Data = []byte("Tunnel added\n")
	log.Println("Tunnel added.")
	return result, nil
}

func (h *handler) removeTunnel(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToTunnel(packet.Data)
	if err != nil {
		log.Println("Uanble to decode tunnel")
		result.Data = []byte("Uanble to decode tunnel\n")
		return result, err
	}
	err = h.configRepo.RemoveTunnel(doc.Name)
	if err != nil {
		log.Println("failed to remove tunnel.", err)
	}
	result.Data = []byte("Tunnel removed\n")
	log.Println("Tunnel removed.")
	return result, nil
}

func (h *handler) listTunnel(_ models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	docs, err := h.configRepo.ListTunnel()
	if err != nil {
		log.Println("failed to read tunnels.", err)
	}

	liveConns := h.socketRepo.LiveConnections()

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Name", "Status", "Local Port", "Remote Host", "Remote Port", "Host Chain"})

	for _, doc := range docs {
		status := "Ready"
		liveConn, ok := liveConns[doc.Name]
		if ok {
			if liveConn.IsActive() {
				status = "Running"
			}
		}
		table.Append([]string{doc.Name, status, strconv.Itoa(doc.LocalPort), doc.RemoteHost, strconv.Itoa(doc.RemotePort), strings.Join(doc.HostChain, ",")})
	}

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	log.Println("TUnnel removed.")
	return result, nil
}

func (h *handler) detailTunnel(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   "string",
		Data:   []byte("\n"),
	}
	doc, err := utils.ToTunnel(packet.Data)
	if err != nil {
		log.Println("Uanble to decode tunnel")
		result.Data = []byte("Uanble to decode tunnel\n")
		return result, err
	}

	doc, err = h.configRepo.DetailTunnel(doc.Name)
	if err != nil {
		log.Println("failed to read tunnel.", err)
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Key", "Value"})

	table.Append([]string{"Name:", doc.Name})
	table.Append([]string{"Local Port:", strconv.Itoa(doc.LocalPort)})
	table.Append([]string{"Remote Host:", doc.RemoteHost})
	table.Append([]string{"Remote Port:", strconv.Itoa(doc.RemotePort)})
	table.Append([]string{"Host Chain :", strings.Join(doc.HostChain, ",")})

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	return result, nil
}

func (h *handler) startTunnel(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}
	doc, err := utils.ToTunnel(packet.Data)
	if err != nil {
		log.Println("Uanble to decode tunnel")
		result.Data = []byte("Uanble to decode tunnel\n")
		return result, err
	}

	doc, err = h.configRepo.DetailTunnel(doc.Name)
	if err != nil {
		log.Println("failed to read tunnel.", err)
	}

	go h.socketRepo.Start(doc)

	result.Data = []byte("Tunnel started successfully.\n")
	return result, nil
}

func (h *handler) stopTunnel(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}
	doc, err := utils.ToTunnel(packet.Data)
	if err != nil {
		log.Println("Uanble to decode tunnel")
		result.Data = []byte("Uanble to decode tunnel\n")
		return result, err
	}

	doc, err = h.configRepo.DetailTunnel(doc.Name)
	if err != nil {
		log.Println("failed to read tunnel.", err)
	}

	h.socketRepo.Stop(doc)

	result.Data = []byte("Tunnel stopped successfully.\n")
	return result, nil
}
