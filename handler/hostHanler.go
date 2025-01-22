package handler

import (
	"bytes"
	"log"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/utils"
)

func (h *handler) addHost(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToHostConfig(packet.Data)
	if err != nil {
		log.Println("Uanble to decode host")
		return packet, err
	}
	err = h.configRepo.AddHost(doc.Name, doc.Host, doc.Port, doc.User, doc.CredientialName)
	if err != nil {
		log.Println("failed to add host.", err)
	}
	result.Data = []byte("Host added\n")
	log.Println("Host added.")
	return result, nil
}

func (h *handler) removeHost(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToHostConfig(packet.Data)
	if err != nil {
		log.Println("Uanble to decode host")
		result.Data = []byte("Uanble to decode host\n")
		return result, err
	}

	if len(doc.Name) == 0 {
		result.Data = []byte("Host name is required\n")
		return result, nil
	}

	err = h.configRepo.RemoveHost(doc.Name)
	if err != nil {
		log.Println("failed to remove host.", err)
	}
	result.Data = []byte("Host removed\n")
	log.Println("Host removed.")
	return result, nil
}

func (h *handler) listHost(_ models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	docs, err := h.configRepo.ListHost()
	if err != nil {
		log.Println("failed to read hosts.", err)
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Name", "Host", "Port", "User", "Crediential Name"})

	for _, doc := range docs {
		table.Append([]string{doc.Name, doc.Host, strconv.Itoa(doc.Port), doc.User, doc.CredientialName})
	}

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	return result, nil
}

func (h *handler) detailHost(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}
	doc, err := utils.ToHostConfig(packet.Data)
	if err != nil {
		log.Println("Uanble to decode host")
		result.Data = []byte("Uanble to decode host\n")
		return result, err
	}

	if len(doc.Name) == 0 {
		result.Data = []byte("Host name is required\n")
		return result, nil
	}

	doc, err = h.configRepo.DetailHost(doc.Name)
	if err != nil {
		log.Println("failed to read host.", err)
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Key", "Value"})

	table.Append([]string{"Name:", doc.Name})
	table.Append([]string{"Host:", doc.Host})
	table.Append([]string{"Port:", strconv.Itoa(doc.Port)})
	table.Append([]string{"User :", doc.User})
	table.Append([]string{"Crediential Name :", doc.CredientialName})

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	return result, nil
}
