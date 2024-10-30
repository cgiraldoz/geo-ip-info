package geolite

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

type GeoLite struct{}

func NewGeoLite() *GeoLite {
	return &GeoLite{}
}

func (g *GeoLite) GetLocation(ip string) *geoip2.City {

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}

	defer func(db *geoip2.Reader) {
		closingError := db.Close()
		if closingError != nil {
			panic(closingError)
		}
	}(db)

	record, err := db.City(net.ParseIP(ip))

	if err != nil {
		panic(err)
	}

	return record
}
