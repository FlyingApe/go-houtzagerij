package main

var latten []lat // slice met alle gezaagde latten

type lat struct {
	Restant    int
	UniekLatId int
	//	Zaaglengten []zaaglengte
}

var UniekId int = 0

func nieuweLat() lat {
	nieuweLat := lat{300, UniekId}
	UniekId++
	return nieuweLat
}

func (l *lat) Zaag(z zaaglengte) {
	l.Restant = l.Restant - z.Lengte
}

func (l *lat) Herzaag() {
	l.Restant = 300
	for _, v := range alleLengten.lengten {
		if v.LatId == l.UniekLatId {
			l.Restant = l.Restant - v.Lengte
		}
	}
}

func getLatByUniekLatId(UniekLatId int) lat {
	var returnLat lat
	for _, v := range latten {
		if v.UniekLatId == UniekLatId {
			returnLat = v
		}
	}
	return returnLat
}
