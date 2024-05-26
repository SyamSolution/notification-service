package model

type DataMessage struct {
	OrderID      string  `json:"order_id"`
	Email        string  `json:"email"`
	URL          string  `json:"url"`
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	DeadlineDate string  `json:"deadline_date"`
	Total        float32 `json:"total"`
}

type EmailPDFMessage struct {
	Email          string         `json:"email"`
	OrderId        string         `json:"order_id"`
	EventName      string         `json:"event_name"`
	Price          float32        `json:"price"`
	NumberOfTicket int            `json:"number_of_ticket"`
	EventDate      string         `json:"event_date"`
	EventTime      string         `json:"event_time"`
	Venue          string         `json:"venue"`
	CustomerName   string         `json:"customer_name"`
	PurchaseDate   string         `json:"purchase_date"`
	DetailTickets  []DetailTicket `json:"detail_tickets"`
}

type DetailTicket struct {
	TicketType  string `json:"ticket_type"`
	TotalTicket int    `json:"total_ticket"`
}
