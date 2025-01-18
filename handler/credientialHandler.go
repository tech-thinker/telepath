package handler

import (
	"bytes"
	"log"

	"github.com/olekukonko/tablewriter"
	"github.com/tech-thinker/telepath/constants"
	"github.com/tech-thinker/telepath/models"
	"github.com/tech-thinker/telepath/utils"
)

func (h *handler) addCrediential(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToCrediential(packet.Data)
	if err != nil {
		log.Println("Uanble to decode crediential")
		return packet, err
	}
	err = h.configRepo.AddCrediential(doc.Name, doc.Type, doc.Password, doc.KeyFile)
	if err != nil {
		log.Println("failed to add crediential.", err)
	}
	result.Data = []byte("crediential added\n")
	log.Println("Crediential added.")
	return result, nil
}

func (h *handler) removeCrediential(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	doc, err := utils.ToCrediential(packet.Data)
	if err != nil {
		log.Println("Uanble to decode crediential")
		result.Data = []byte("Uanble to decode crediential\n")
		return result, err
	}
	err = h.configRepo.RemoveCrediential(doc.Name)
	if err != nil {
		log.Println("failed to remove crediential.", err)
	}
	result.Data = []byte("crediential removed\n")
	log.Println("Crediential removed.")
	return result, nil
}

func (h *handler) listCrediential(_ models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}

	docs, err := h.configRepo.ListCrediential()
	if err != nil {
		log.Println("failed to remove crediential.", err)
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Name", "Type"})

	for _, doc := range docs {
		table.Append([]string{doc.Name, doc.Type})
	}

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	log.Println("Crediential removed.")
	return result, nil
}

func (h *handler) detailCrediential(packet models.Packet) (models.Packet, error) {
	result := models.Packet{
		Action: "result",
		Type:   constants.PACKET_TYPE_STRING,
		Data:   []byte("\n"),
	}
	cred, err := utils.ToCrediential(packet.Data)
	if err != nil {
		log.Println("Uanble to decode crediential")
		result.Data = []byte("Uanble to decode crediential\n")
		return result, err
	}

	cred, err = h.configRepo.DetailCrediential(cred.Name)
	if err != nil {
		log.Println("failed to read crediential.", err)
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Key", "Value"})

	table.Append([]string{"Name:", cred.Name})
	table.Append([]string{"Type:", cred.Type})
	table.Append([]string{"Password:", cred.Password})
	table.Append([]string{"Key File:", cred.KeyFile})

	table.SetBorder(true) // Enable/Disable borders
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	result.Data = buf.Bytes()
	return result, nil
}
