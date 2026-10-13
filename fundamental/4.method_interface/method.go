package main

import "fmt"

type Album struct {
	title        string
	release_date string
	genre        string
	sellCount    int
}
type Singer struct {
	name    string
	country string
	age     int
	albums  []Album
}

func (s Singer) getBestSellerAlbum() Album {
	max, position := 0, 0
	for index, v := range s.albums {
		if v.sellCount > max {
			max = v.sellCount
			position = index
		}
	}
	return s.albums[position]
}
func (s *Singer) setNameSinger(name string) *Singer {
	s.name = name
	return s
}
func valueReceiver() {
	var Katy = Singer{
		name:    "Katy perry",
		country: "American",
		age:     23,
		albums: []Album{
			{
				title:        "Album1",
				release_date: "04/05/2004",
				genre:        "344",
				sellCount:    32,
			},
			{
				title:        "Album2",
				release_date: "03/05/2004",
				genre:        "23",
				sellCount:    39,
			},
		},
	}
	maxAlbum := Katy.getBestSellerAlbum()
	fmt.Println("Album bán chạy nhất là", maxAlbum.title)
	newKaty := Katy.setNameSinger("Harry Styles")
	fmt.Println("New singer", newKaty.name)

}
func Methods() {
	valueReceiver()
}
