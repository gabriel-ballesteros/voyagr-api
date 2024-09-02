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
	EventDatetime string `bson:"eventDatetime"`
	Notes         string `bson:"notes"`
}

type Trip struct {
	ID          string             `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Start       string             `bson:"start"`
	End         string             `bson:"end"`
	Owner       string             `bson:"owner"`
	SharedWith  []string           `bson:"sharedWith"`
	Itinerary   []ItineraryElement `bson:"itinerary"`
}
