package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string               `bson:"name" json:"name"`
	Address string               `bson:"address" json:"address"`
	Rooms   []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating  int                  `bson:"rating" json:"rating"`
}

type RoomType int

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size" json:"size"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
