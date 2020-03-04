package main

import (
	"fmt"
	"math/rand"
	"time"

	goclear "github.com/TwinProduction/go-clear"
	"github.com/segmentio/go-prompt"
)

//deck global var
var deck []string

// action pioche
var promptPioche = []string{
	"Piocher dans la pioche", "Piocher dans la defausse", "Dire Gabo",
}

// action apres pioche
var promptCard = []string{
	"Remplacer une carte de son jeu", "Jeter", "Utiliser pouvoir", "Defausser un carte en double",
}

var mainReplace []string

// main des joeurs et var
var mainJ1, mainJ2 []string
var grave, card string
var win = 0

func main() {
	goclear.Clear()

	// debut partie
	initPartie()
	// lancement partie
	partie()

}
func initPartie() {
	//init du deck
	deck = append(deck, "1co", "2co", "3co", "4co", "5co", "6co", "7co", "8co", "9co", "Xco", "Vco", "Dco", "Rco")
	deck = append(deck, "1ca", "2ca", "3ca", "4ca", "5ca", "6ca", "7ca", "8ca", "9ca", "Xca", "Vca", "Dca", "Rca")
	deck = append(deck, "1tr", "2tr", "3tr", "4tr", "5tr", "6tr", "7tr", "8tr", "9tr", "Xtr", "Vtr", "Dtr", "Rtr")
	deck = append(deck, "1pi", "2pi", "3pi", "4pi", "5pi", "6pi", "7pi", "8pi", "9pi", "Xpi", "Vpi", "Dpi", "Rpi")

	// shuffle hand
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	// init grave
	grave = deck[len(deck)-1]
	deck = deck[:len(deck)-1]

	//init des mains
	mainJ1 = initMain(mainJ1, 4)
	mainJ2 = initMain(mainJ2, 4)

}

func initMain(mainJx []string, n int) []string {

	// loop selon nb de card
	for i := 0; i < n; i++ {
		// deplacement d'une carte du deck
		card := deck[len(deck)-1]
		deck = deck[:len(deck)-1]
		//ajout de la carte a la main
		mainJx = append(mainJx, card)
		//fmt.Println("mainJx", mainJx)

	}
	return mainJx

}

func funcTour(mainJx []string, computer bool) []string {
	//0"Piocher dans la pioche", 1"Piocher dans la defausse", 2"Dire Gabo"
	i := prompt.Choose("Quelle action ?", promptPioche)

	//fmt.Println("picked ", promptPioche[i])

	printAll()

	switch i {
	case 0:
		card = funcPioche("pioche")
		fmt.Println("Carte recuperée dans la pioche ", card)
	case 1:
		card = funcPioche("grave")
		fmt.Println("Carte recuperée dans la defausse ", card)
	case 2:
		fmt.Println("GABOOOOOOOOO")
	}

	if i != 2 { // if not gabo
		//0"Remplacer une carte de son jeu", 1"Jeter", 2"Utiliser pouvoir", 3"Defausser un carte en double"
		y := prompt.Choose("Que faire avec cette carte ?", promptCard)
		printAll()
		//fmt.Println("picked ", promptCard[y])
		switch y {
		case 0:
			mainJx = printMainhide(mainJx, card, "replace")
		case 1:
			//fmt.Println("Defausse de la carte ", card)
			grave = card
		case 2:
			// use power
		case 3:
			//deffauss double
			mainJx = printMainhide(mainJx, card, "double")
		}
	} else {
		//gabo
		gabo(computer)
	}

	return mainJx

}

func gaboCount(mainJx []string) int {
	var count = 0
	for i := range mainJx {
		switch mainJx[i][0] {
		case 'R':
			if mainJx[i] == "Rca" || mainJx[i] == "Rco" {
				count += 0
			} else {
				count += 15
			}
		case 'X', 'V', 'D':
			count += 10
		//	fmt.Println("countJ1 ", countJ1)
		default:
			// - 48 pour convertir en int
			count += int(mainJx[i][0] - 48)
			//	fmt.Println("countJ1 ", countJ1)

		}
	}
	return count
}

func gabo(computer bool) {
	goclear.Clear()

	// calcul total de point
	countJ1 := gaboCount(mainJ1)
	countJ2 := gaboCount(mainJ2)

	// calcul du gagnant
	var msgGagant string
	if countJ1 == countJ2 {
		msgGagant = "Equalite parfaite !"
	} else if computer == false && countJ1 < countJ2 && countJ1 < 7 {
		// si c'est le joueur 1 qui a dit gabo
		msgGagant = "Le joueur 1 gagne !"
	} else {
		msgGagant = "Le joueur 2 gagne !"
	}

	fmt.Println("---------------PARTIE TERMINEE----------------- ")
	fmt.Println("X,V,D = 10pts | Rca,Rco = 0pts | Rpi, Rtr = 15pts")
	fmt.Println("------------------------------------------------ ")
	fmt.Println("mainJ1 ", mainJ1)
	fmt.Println("countJ1 ", countJ1)
	fmt.Println("mainJ2 ", mainJ2)
	fmt.Println("countJ2 ", countJ2)
	fmt.Println("------------------------------------------------ ")
	fmt.Println(msgGagant)
	fmt.Println("------------------------------------------------ ")

	// end partie
	win = 1

}

func funcPioche(zone string) string {

	if zone == "pioche" {
		card = deck[len(deck)-1]
		deck = deck[:len(deck)-1]
	} else { // grave
		card = grave
	}

	return card

}

func printMainhide(mainJx []string, cardAction string, action string) []string {
	//affiche la main en masquée pour remplacement ou double
	var mainReplace []string
	for i := range mainJx {
		i = i + 1
		mainReplace = append(mainReplace, "X")
	}

	if action == "replace" {
		// replace
		z := prompt.Choose("Quelle carte remplacer ?", mainReplace)
		grave = mainJx[z]
		mainJx[z] = cardAction
	} else {
		// double
		z := prompt.Choose("Quelle carte est en double ?", mainReplace)

		if mainJx[z][0] == cardAction[0] {
			grave = mainJx[z]
			mainJx = removeSlice(mainJx, z)
		} else {
			// ce n'est pas un double
			mainJx = append(mainJx, cardAction)
		}

	}

	return mainJx

}

func removeSlice(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func printAll() {
	goclear.Clear()

	fmt.Println("---------------TEST----------------- ")
	fmt.Println("mainJ1 ", mainJ1)
	fmt.Println("mainJ2 ", mainJ2)
	fmt.Println("grave ", grave)
	fmt.Println("deck ", deck)
	fmt.Println("---------------TEST----------------- ")

}

func partie() {

	for win != 1 {

		printAll()

		// bool pour indiquer si computer ou no
		mainJ1 = funcTour(mainJ1, false)
	}

}
