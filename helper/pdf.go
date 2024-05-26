package helper

import (
	"bytes"
	"encoding/base64"
	"log"
	"strconv"

	"github.com/SyamSolution/notification-service/internal/model"
	"github.com/signintech/gopdf"
)

func GeneratePDF(message model.EmailPDFMessage) (string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	
	err := pdf.Image("./assets/image/header.png", 0, 0, nil)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}
	err = pdf.AddTTFFont("times", "./assets/font/times.ttf")
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	err = pdf.SetFont("times", "", 14)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	leftMargin := 40.0
	topMargin := 40.0
	pdf.SetLeftMargin(leftMargin)
	pdf.SetTopMargin(topMargin)

	pageWidth := gopdf.PageSizeA4.W - leftMargin

	textWidth, err := pdf.MeasureTextWidth("CONCERT MUSIC 2024")
	if err != nil {
		log.Print(err.Error())
	}
	x := (pageWidth - textWidth) / 2 + 20
	y := 170.0
	
	pdf.SetX(x)
	pdf.SetY(y)

	pdf.SetFont("times", "B", 30)
	pdf.Cell(nil, "CONCERT MUSIC 2024")

	pdf.SetFont("times", "", 14)
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
	pdf.WritePdf("ticket.pdf")

	// Convert the byte buffer to a base64 string
	pdfBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return pdfBase64, nil
}
