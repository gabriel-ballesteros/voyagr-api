package domain

type ItineraryElement struct {
	Title         string `json:"title"`
	Type          string `json:"type"`
	From          string `json:"from"`
	To            string `json:"to"`
	Departure     string `json:"departure"`
	Arrival       string `json:"arrival"`
	Address       string `json:"address"`
	FlightStatus  string `json:"status"`
	FlightGate    string `json:"gate"`
	Seat          string `json:"seat"`
	PaymentStatus string `json:"paymentStatus"`
	CheckIn       string `json:"checkIn"`
	CheckOut      string `json:"checkOut"`
	Notes         string `json:"notes"`
}

type Trip struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Start      string             `json:"start"`
	End        string             `json:"end"`
	Owner      int                `json:"owner"`
	SharedWith []int              `json:"sharedWith"`
	Itinerary  []ItineraryElement `json:"itinerary"`
}
