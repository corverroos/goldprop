package models

import "github.com/jinzhu/gorm"

type ListingType int

const (
	ListingTypeUnknown ListingType = 0
	ListingTypeSale    ListingType = 1
	ListingTypeRent    ListingType = 2
)

type PropertyType int

const (
	PropertyTypeUnknown   PropertyType = 0
	PropertyTypeApartment PropertyType = 1
	PropertyTypeHouse     PropertyType = 2
)

type Site int

const (
	SiteUnknown Site = 0
	SiteP24     Site = 1
)

type Listing struct {
	gorm.Model

	Listing      ListingType
	PropertyType PropertyType
	Site         Site
	SiteID       string
	Location     string
	Address      string
	Price        int    // in Rand
	Bedrooms     int
	Bathrooms    int
	Garages      int
	FloorSize    int // m^2
	ErfSize      int // m^2
	URL          string `sql:"size:4096",gorm:"unique;not null"`
	Description  string `sql:"type:text"`
}

type Features struct {
	gorm.Model

	ListingID uint
	Key       string
	Value     string
}
