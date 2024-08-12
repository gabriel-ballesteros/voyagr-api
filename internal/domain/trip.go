package domain

type PlaneInfo struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Address   string `json:"address"`
	Status    string `json:"status"`
	Gate      string `json:"gate"`
	Seat      string `json:"seat"`
	Notes     string `json:"notes"`
}
type TrainBusInfo struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	From      string `json:"from"`
	To        string `json:"to"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Address   string `json:"address"`
	Platform  string `json:"platform"`
	Seat      string `json:"seat"`
	Notes     string `json:"notes"`
}
type CarInfo struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	From  string `json:"from"`
	To    string `json:"to"`
	Notes string `json:"notes"`
}
type AccomodationInfo struct {
	Title         string `json:"title"`
	Type          string `json:"type"`
	Address       string `json:"address"`
	PaymentStatus string `json:"paymentStatus"`
	CheckIn       string `json:"checkIn"`
	CheckOut      string `json:"checkOut"`
	Notes         string `json:"notes"`
}
type EventInfo struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Notes   string `json:"notes"`
}

type Trip struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Start      string        `json:"start"`
	End        string        `json:"end"`
	Owner      int           `json:"owner"`
	SharedWith []int         `json:"sharedWith"`
	Itinerary  []interface{} `json:"itinerary"`
}
