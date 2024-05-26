package helper

import (
	"bytes"
	"encoding/base64"
	"github.com/SyamSolution/notification-service/internal/model"
	"github.com/signintech/gopdf"
	"log"
	"strconv"
)

func GeneratePDF(message model.EmailPDFMessage) (string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("times", "./times.ttf")
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	err = pdf.SetFont("times", "", 14)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	pdf.Cell(nil, "Concert Music 2024")
	pdf.Br(40)
	pdf.Cell(nil, "Ticket Order")
	pdf.Br(40)
	pdf.Cell(nil, "Order ID: "+message.OrderId)
	pdf.Br(20)
	pdf.Cell(nil, "Event Name: "+message.EventName)
	pdf.Br(20)
	pdf.Cell(nil, "Price: "+strconv.Itoa(int(message.Price)))
	pdf.Br(20)
	pdf.Cell(nil, "Number of Ticket: "+strconv.Itoa(message.NumberOfTicket))
	pdf.Br(20)
	pdf.Cell(nil, "Event Date: "+message.EventDate)
	pdf.Br(20)
	pdf.Cell(nil, "Event Time: "+message.EventTime)
	pdf.Br(20)
	pdf.Cell(nil, "Venue: "+message.Venue)
	pdf.Br(20)
	pdf.Cell(nil, "Customer Name: "+message.CustomerName)
	pdf.Br(20)
	pdf.Cell(nil, "Purchase Date: "+message.PurchaseDate)
	pdf.Br(20)
	pdf.Cell(nil, "Detail Ticket: ")
	pdf.Br(20)
	for i, ticket := range message.DetailTickets {
		pdf.Cell(nil, strconv.Itoa(i+1)+". Ticket type: "+ticket.TicketType)
		pdf.Br(20)
		pdf.Cell(nil, "    Total ticket: "+strconv.Itoa(ticket.TotalTicket))
		pdf.Br(20)
	}

	var buf bytes.Buffer
	err = pdf.Write(&buf)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	// Convert the byte buffer to a base64 string
	pdfBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return pdfBase64, nil
}
