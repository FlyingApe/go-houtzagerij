package main

import (
	"math/rand"
	"sort"
	"time"
)

/*
Zaaglengte
*/
type zaaglengte struct {
	Lengte, LatId int
}

func (z *zaaglengte) AssignLat(latId int) {
	z.LatId = latId
}

/*
Zaaglengten
*/
type zaaglengten struct {
	lengten []zaaglengte
}

func (z *zaaglengten) AddLengte(lengte zaaglengte) {
	z.lengten = append(z.lengten, lengte)
	z.SortHighLow()
}

func (z *zaaglengten) fillRandom(aantal int) {
	j := 0
	for i := 0; i < aantal; i++ {
		rand.Seed(int64(time.Now().Nanosecond() + i))
		randInt := rand.Intn(59) + 1
		lengte := zaaglengte{randInt * 5, -1}
		if j == 2 || j == 4 {
			for lengte.Lengte > 130 {
				rand.Seed(int64(time.Now().Nanosecond() + i))
				randInt = rand.Intn(59) + 1
				lengte = zaaglengte{randInt * 5, -1}
			}
			if j == 4 {
				j = 0
			}
		}
		if lengte.Lengte != 300 {
			z.AddLengte(lengte)
		}
		j++
	}
}

func (z *zlComb) AddLengte(lengte zaaglengte, UId int) {
	uniekeLengte := uniekeZaaglengte{lengte, UId}
	z.Zaaglengten = append(z.Zaaglengten, uniekeLengte)
	//z.SortHighLow()
}

func (z *zaaglengten) ReturnZlCombByLatId(latId int) (returnVar zlComb) {
	UnieklatId := latten[latId].UniekLatId
	for i, v := range z.lengten {
		if v.LatId == UnieklatId {
			returnVar.AddLengte(v, i)
		}
	}
	return
}

func (z *zaaglengten) ReturnByLatId2(latId int) (returnVar zaaglengten) {
	UnieklatId := latten[latId].UniekLatId
	for _, v := range z.lengten {
		if v.LatId == UnieklatId {
			returnVar.AddLengte(v)
		}
	}
	return
}

func (z *zaaglengten) SortHighLow() {
	sort.SliceStable(z.lengten, func(i, j int) bool { return z.lengten[i].Lengte > z.lengten[j].Lengte })
}

/*
Combinaties
*/
type uniekeZaaglengte struct {
	zaaglengte
	uniekLengteId int // positie binnen de oorspronkelijke lat
}

type zlComb struct { // combinatie van zaaglengtes
	Zaaglengten []uniekeZaaglengte
	Rest        int
	//	Minlengte   int
	//	Maxlengte   int
}

type wissel struct {
	combId, lengteId, grootsteRest int
}
