package main

import (
	"fmt"
	"os"
	"os/exec"
)

var win,draw,loose bool = false, false, false
var player1, player2 string = "", ""
var player1Active = true
var winner , looser = "" , ""

func getWinner(f int) {
	if f == 1 {
		winner = player1
		looser = player2
	} else {
		winner = player2
		looser = player1
	}
}

func checkState(a, b, c int, state *[]int) {
	if (*state)[a] == (*state)[b] && (*state)[b] == (*state)[c] && ((*state)[a] == 1 || (*state)[a] == 2) {
		win = true
		getWinner((*state)[a])
	}
}

func GetResult(winningPatterns *[][]int, state *[]int) {
	for i := 0; i < len(*winningPatterns); i++ {
		checkState((*winningPatterns)[i][0], (*winningPatterns)[i][1], (*winningPatterns)[i][2], &(*state))
		if win {
			break
		}
	}
	chekcIfExhausted(&(*state))
}

func chekcIfExhausted(state *[]int) {
	exhausted := true
	for i := 0; i < len(*state); i++ {
		if (*state)[i] == 0 {
			exhausted = false
			break
		}
	}
	if exhausted && !win && !loose {
		draw = true
	}
}
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string  `json:"error"`
}

func ValiDate(val int) {
	if val == 1 {
		GetValidate(1)
	}
	if val == 2 {
		GetValidate(2)
	}
}

func GetValidate(val int) {
	if val == 1  {
		fmt.Println("Please enter player 1 name :")
		fmt.Scanln(&player1)
		if player1 == "" {
			ValiDate(1)
		}else {
			return
		}
	}

	if val == 2  {
		fmt.Println("Please enter player 2 name :")
		fmt.Scanln(&player2)
		if player2 == "" {
			ValiDate(2)
		}else {
			return
		}
	}
}

func GetNames() {
	fmt.Println("Please enter player 1 name :")
	fmt.Scanln(&player1)
	if player1 == "" {
		ValiDate(1)
	}
	fmt.Println("Please enter player 2 name :")
	fmt.Scanln(&player2)
	if player2 == "" {
		ValiDate(2)
	}

}

func printState(state *[]int) {
	var toPrintArr [9]string
	fmt.Println(state)
	for i := 0; i < len(*state); i++ {
		if (*state)[i] == 0 {
			toPrintArr[i] = "="
		} else if (*state)[i] == 1 {
			toPrintArr[i] = "X"
		} else if (*state)[i] == 2 {
			toPrintArr[i] = "O"
		}
	}

	fmt.Printf("%s   %s   %s \n", toPrintArr[0], toPrintArr[1], toPrintArr[2])
	fmt.Printf("%s   %s   %s \n", toPrintArr[3], toPrintArr[4], toPrintArr[5])
	fmt.Printf("%s   %s   %s \n", toPrintArr[6], toPrintArr[7], toPrintArr[8])
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func start(state *[]int, winningPatterns *[][]int) {
	GetNames()
	for !win && !draw && !loose {
		clearScreen()
		printState(&(*state))
		var input int
		activePlayerName := ""
		activePlayerVal := 1

		if player1Active {
			activePlayerName = player1
			activePlayerVal = 1
		} else {
			activePlayerName = player2
			activePlayerVal = 2
		}

		for {
			fmt.Printf("%s's turn now \n", activePlayerName)
			_, err := fmt.Scanln(&input)
			if err != nil {
				var discard string
				fmt.Scanln(&discard)
				clearScreen()
				printState(&(*state))
				continue
			}
			fmt.Println(activePlayerVal)
			fmt.Println("visa")

			if input != 0 && input >= 1 && input <= 9 {
				input--
				if (*state)[input] == 0 {
					(*state)[input] = activePlayerVal // allo to modify the state
					break
				}
			}
			fmt.Println((*state)[input])
			fmt.Println(activePlayerVal)

			clearScreen()
			printState(&(*state))
		}
		GetResult(&(*winningPatterns), &(*state))
		player1Active = !player1Active
	}
	printState(&(*state))
	if win {
		fmt.Printf("\n %s Winner party to banti he!!!! 😀 😀 😀 😀 😀", winner+" => ")
		fmt.Printf("\n %s Looser 😒 😒 😒 😒 😞", looser+" => ")
	}
	if draw {
		fmt.Println("Match Draw")
	}
}

func main() {
	state := []int{0, 0, 0, 0, 0, 0, 0, 0, 0} // initialize the state
	winningPatterns := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}
	start(&state, &winningPatterns)
}
