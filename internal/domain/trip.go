package domain

type ItineraryElement struct {
	Title         string `bson:"title"`
	Type          string `bson:"type"`
	From          string `bson:"from"`
	To            string `bson:"to"`
	Departure     string `bson:"departure"`
	Arrival       string `bson:"arrival"`
	Address       string `bson:"address"`
	FlightStatus  string `bson:"status"`
	FlightGate    string `bson:"gate"`
	Seat          string `bson:"seat"`
	PaymentStatus string `bson:"paymentStatus"`
	CheckIn       string `bson:"checkIn"`
	CheckOut      string `bson:"checkOut"`
	Notes         string `bson:"notes"`
}

type Trip struct {
	ID          string             `bson:"_id,omitempty"`
	Description string             `bsono:description`
	Name        string             `bson:"name"`
	Start       string             `bson:"start"`
	End         string             `bson:"end"`
	Owner       string             `bson:"owner"`
	SharedWith  []int              `bson:"sharedWith"`
	Itinerary   []ItineraryElement `bson:"itinerary"`
}
