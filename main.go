package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gonum/stat/combin"
)

var alleLengten zaaglengten

func main() {
	var run bool = true
	for run {

		if len(alleLengten.lengten) != 0 {
			fmt.Printf("\nOpgeslagen zaaglengtes:")
			for i := 0; i < len(alleLengten.lengten); i++ {
				fmt.Printf(" %d", alleLengten.lengten[i].Lengte)
			}
			fmt.Printf("\n\n")
		}

		maat := scan()
		switch maat {
		case 600:
			continue
		case 400:
			run = false
			continue
		case 500:
			assignLengten()
			fmt.Printf("Verdeling voor optimalisatie:\n\n")
			printCombinaties()

			// grove toewijzing is nu gedaan, nu kunnen we kijken of er nog efficiëntiewinst te halen is
			returnValue := hercombineerLatten()
			_ = returnValue
			fmt.Printf("\nVerdeling na optimalisatie:\n\n")
			printCombinaties()

			// programma variabelen terug naar beginwaarden zetten
			alleLengten.lengten = nil
			latten = nil
			UniekId = 0
			continue
		default:
			l := zaaglengte{maat, -1}
			alleLengten.AddLengte(l)
		}
	}
}

func scan() int {
	var return_var int

	//scanner aanmaken
	scanner := bufio.NewScanner(os.Stdin)

	//vraag, scan en print getal
	var end_loop bool = false
	for !end_loop {
		fmt.Printf("--- h voor help ---\n")
		scanner.Scan()
		input := scanner.Text()
		switch input {
		case "h":
			fmt.Printf("--- -----------------------------------------------------------------------------+\n")
			fmt.Printf("--- \n")
			fmt.Printf("--- Geef een maat in cm tussen 0 en 300 en druk op enter \n")
			fmt.Printf("--- \n")
			fmt.Printf("--- 'w'		- voor willekeurige maten. \n")
			fmt.Printf("--- 'a'		- voor de invoer van een volledige array van maten \n")
			fmt.Printf("--- 'Enter'	- voor bepaling optimaal zaagplan \n")
			fmt.Printf("--- \n")
			fmt.Printf("--- 'q' 	- voor exit \n")
			fmt.Printf("--- \n")
			fmt.Printf("--- -----------------------------------------------------------------------------+\n\n")
			return_var = 600
			end_loop = true
		case "w":
			fmt.Printf("Aantal willikeurige nummers dat toegevoegd moet worden: ")
			scanner.Scan()
			input := scanner.Text()

			aantal, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("fout: ongeldige invoer")
			} else {
				alleLengten.fillRandom(aantal)
				return_var = 600
				end_loop = true
			}
		case "a":
			fmt.Printf("Array die ingevoerd moet worden: ")
			scanner.Scan()
			input := scanner.Text()
			waarden := strings.Split(input, " ")

			for _, v := range waarden {
				w, err := strconv.Atoi(v)
				if err != nil {
					fmt.Println(err)
				} else {
					lengte := zaaglengte{w, -1}
					alleLengten.AddLengte(lengte)
					return_var = 600
					end_loop = true
				}
			}
		case "q":
			return_var = 400
			end_loop = true
		case "":
			return_var = 500
			end_loop = true
		default:
			number, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("fout: ongeldige invoer")
			} else if number == 0 || number >= 301 {
				fmt.Println("fout: ongeldige maat")
			} else {
				return_var = number
				end_loop = true
			}
		}
	}

	return return_var
}

func koppelLengteEnLat(latId, lengteId int) {
	for i, _ := range latten {
		if latten[i].UniekLatId == latId {
			latten[i].Zaag(alleLengten.lengten[lengteId])
			alleLengten.lengten[lengteId].AssignLat(latten[i].UniekLatId)
		}
	}
}

func assignLengten() {
	latten = append(latten, nieuweLat()) // slice met alle gezaagde latten bevat nu 1 ongezaagde lat

	// dit algoritme wijst de grootste nog aanwezige zaaglengte toe aan de kleinst aanwezige restlengte (of een nieuwe lat)
	m := len(alleLengten.lengten)
	for lengteId := 0; lengteId < m; lengteId++ { // voor elke ingevoerde zaaglengte
		sort.SliceStable(latten, func(k, l int) bool { return latten[k].Restant < latten[l].Restant })
		n := len(latten)
		for latId := 0; latId < n; latId++ { // voor elke bestaande lat
			if latten[latId].Restant >= alleLengten.lengten[lengteId].Lengte { // als de zaaglengte in het restant van deze lat past
				koppelLengteEnLat(latten[latId].UniekLatId, lengteId)
				break
			} else if latId == n-1 && latten[latId].Restant < alleLengten.lengten[lengteId].Lengte { // als dit de laatste lat is en de zaaglengte past er niet meer in
				latten = append(latten, nieuweLat())
				koppelLengteEnLat(latten[UniekId-1].UniekLatId, lengteId)
				//fmt.Println(n, latId)
				break
			}
		}
	}
}

func hercombineerLatten() bool {
	// controle of herordenen van zaaglengte efficiëntiewinst op kan brengen
	numLatten := len(latten)
	var rest int
	for latId := 0; latId < numLatten; latId++ {
		rest += latten[latId].Restant
	}

	// Als de totale restlengte groter is dan 300 cm dan is er mogelijk een efficiëntere combinatie mogelijk.

	// Als de laatste lat echter groter dan een halve lat is, valt er niets te winnen omdat alle aanwezige latten
	// sowieso gezaagd moeten worden omdat elke een zaaglengte van meer dan 150cm bevat)

	//	laatseLat := numLatten - 1
	//	LaatsteLat := alleLengten.ReturnByLatId(laatseLat)

	if rest > 300 { // zoeken naar optimalisatie heeft alleen zin als de restlengten samen groter dan 1 lat zijn
		// optimalisatie van restlengten
		for latId := 0; latId < len(latten); latId++ {
			// optimalisatie heeft alleen zin als het restant van de lat groter dan 0 is
			if latten[latId].Restant != 0 {
				wisselGevonden := true
				for wisselGevonden {
					/* 	verzameling
					//	verzameling Is een collectie met alle mogelijke combinaties van zaaglengtes in deze lat.
					//	Later testen we alle combinaties met elke andere zaaglengte in elke andere lat op efficientiewinst
					*/
					var verzameling []zlComb

					Lengten := alleLengten.ReturnZlCombByLatId(latId)
					numLengten := len(Lengten.Zaaglengten)

					for i, _ := range Lengten.Zaaglengten {
						elements := i + 1
						if numLengten > 2 && elements != numLengten {
							combGen := combin.NewCombinationGenerator(numLengten, elements)
							for combGen.Next() != false {
								comb := combGen.Combination(nil)
								var lengteComb zlComb
								lengteComb.Rest = latten[latId].Restant
								for _, v := range comb {
									lengte := Lengten.Zaaglengten[v]
									lengteComb.Zaaglengten = append(lengteComb.Zaaglengten, lengte)
								}
								verzameling = append(verzameling, lengteComb)
							}
						}
					}

					/* We hebben nu een verzameling van combinaties waarmee we kunnen schuiven en mogelijke efficientiewinst
					// kunnen bereiken. We gaan nu elke combinatie vergelijken met elke zaaglengte en de restlengtes in die latten.
					// Vervolgens kunnen we kijken of er efficientiewinst ontstaat met een wissel. Zoja, dan slaan we deze
					// wissel en de haalbare efficientiewinst tijdelijk op.
					// Tenslotte nemen we de wissel met de hoogste efficientiewinst en voeren deze uit.
					//
					// Een efficientiewinst wordt gedefinieerd als een wissel waarbij een restlengte ontstaat die groter is dan
					// beide oorspronkelijke restlengtes. Uiteindelijk kan op deze manier een restlengte ontstaan waar een
					// extra zaaglengte in past */

					if len(verzameling) != 0 { // er zijn verzamelingen die gewisselt kunnen worden
						//fmt.Printf("lat %d: ", latId+1)
						//fmt.Println(verzameling)

						// eerst kijken welke wissels tussen combinaties van lengten uitgevoerd kunnen worden die grotere restlengtes opleveren
						var wissels []wissel
						for c, v := range verzameling {
							totaalLengte1 := 0
							for _, w := range v.Zaaglengten {
								totaalLengte1 += w.Lengte
							}
							restLengte1 := v.Rest

							for a, w := range alleLengten.lengten {
								if w.LatId != v.Zaaglengten[0].LatId {
									totaalLengte2 := w.Lengte
									restLengte2 := getLatByUniekLatId(w.LatId).Restant
									//evalLengte := w.Lengte+latRest
									if totaalLengte1 > totaalLengte2 && totaalLengte1 <= totaalLengte2+restLengte2 {
										nieuweRest1 := totaalLengte2 + restLengte2 - totaalLengte1
										nieuweRest2 := totaalLengte1 + restLengte1 - totaalLengte2
										if (nieuweRest1 > restLengte1 && nieuweRest1 > restLengte2) || (nieuweRest2 > restLengte1 && nieuweRest2 > restLengte2) {
											// er is optimalisatiewinst mogelijk
											var grootsteRest int
											if nieuweRest1 > nieuweRest2 {
												grootsteRest = nieuweRest1
											} else {
												grootsteRest = nieuweRest2
											}
											wissels = append(wissels, wissel{c, a, grootsteRest})
										}
									}
								}
							}
						}

						// Nu kunnen we de meest efficiente wissel voor deze lat uitvoeren
						if len(wissels) != 0 {
							var besteWissel wissel
							sort.SliceStable(wissels, func(k, l int) bool { return wissels[k].grootsteRest > wissels[l].grootsteRest })
							besteWissel = wissels[0] // deze wissel creërt de grootste restlengte

							// wisselen
							wisselLatId1 := alleLengten.lengten[besteWissel.lengteId].LatId
							wisselLatId2 := -1

							for _, y := range verzameling[besteWissel.combId].Zaaglengten {
								alleLengten.lengten[y.uniekLengteId].LatId = wisselLatId1
								wisselLatId2 = y.LatId
							}

							alleLengten.lengten[besteWissel.lengteId].LatId = wisselLatId2

							// na de wissel de restlengte opnieuw berekenen
							for p, y := range latten {
								if y.UniekLatId == wisselLatId1 || y.UniekLatId == wisselLatId2 {
									latten[p].Herzaag()
								}
							}
						} else { // Er zijn geen nuttige wissel meer mogelijk, de loop voor zoektocht naar efficientie verbreken
							wisselGevonden = false
						}
					} else {
						wisselGevonden = false
					}
				}
			}
		}

		//printCombinaties()

		// lege latten uit de distributie van latten halen
		m := len(alleLengten.lengten)
		for lengteId := m - 1; lengteId >= 0; lengteId-- { // voor elke ingevoerde zaaglengte, van klein naar groot
			n := len(latten)

			sort.SliceStable(latten, func(k, l int) bool { return latten[k].Restant < latten[l].Restant }) // sorteren van de latten op restlengte
			for latId := 0; latId < n; latId++ {                                                           // voor elke bestaande lat
				if latten[latId].Restant == 300 { // als deze lat nu niet meer nodig is
					latten = append(latten[:latId], latten[latId+1:]...)
					continue
				} else if latten[latId].UniekLatId == alleLengten.lengten[lengteId].LatId { // als deze lengte in deze lat zit
					continue
				} else if latten[latId].Restant >= alleLengten.lengten[lengteId].Lengte && alleLengten.lengten[lengteId].LatId == latten[n-1].UniekLatId { // als de zaaglengte in het restant van deze lat past
					//fmt.Printf("LengteId: %d, Lengte: %d \nLat.Restant: %d, ")
					oudLatId := alleLengten.lengten[lengteId].LatId
					koppelLengteEnLat(latten[latId].UniekLatId, lengteId)
					for x, _ := range latten {
						if latten[x].UniekLatId == oudLatId {
							latten[x].Herzaag()
						}
					}
					break
				} else if latId == n-1 { // als dit de laatste lat is en de zaaglengte past er niet meer in
					break
				}
			}
		}
		return true
	} else {
		return false
	}
}

func printCombinaties() {
	var rest int
	fmt.Printf("In totaal moeten er %d latten gezaagd worden. \n", len(latten))
	for i := 0; i < len(latten); i++ {
		rest += latten[i].Restant
		fmt.Printf("Uit lat %d worden de stukken van", i+1)
		z := alleLengten.ReturnByLatId2(i)
		for j := 0; j < len(z.lengten); j++ {
			fmt.Printf(", %d cm", z.lengten[j].Lengte)
		}
		fmt.Printf(" gezaagd. Er blijft %d over.\n", latten[i].Restant)
		//fmt.Printf(" gezaagd. Er blijft %d over. %d \n", latten[i].Restant, z)
	}

	fmt.Printf("Totale restproducten zijn %d cm lang. \n\n", rest)
}
