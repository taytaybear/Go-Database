// starting point for the database application
package main

import (
	"bufio"
	// "fmt"
	"os"

	"go-database/internal/cli"
	"go-database/internal/db"
)

func main() {
	// Display ASCII art banner
	// fmt.Print(`

	//                 ___          ___          ___          ___          ___          ___                            _____
	//     ___        /__/\        /__/\        /  /\        /__/\        /  /\        /  /\        ___               /  /::\      _____
	//    /  /\       \  \:\      |  |::\      /  /:/_      |  |::\      /  /::\      /  /::\      /__/|             /  /:/\:\    /  /::\
	//   /  /:/        \  \:\     |  |:|:\    /  /:/ /\     |  |:|:\    /  /:/\:\    /  /:/\:\    |  |:|            /  /:/  \:\  /  /:/\:\
	//  /__/::\    _____\__\:\  __|__|:|\:\  /  /:/ /:/_  __|__|:|\:\  /  /:/  \:\  /  /:/~/:/    |  |:|           /__/:/ \__\:|/  /:/~/::\
	//  \__\/\:\__/__/::::::::\/__/::::| \:\/__/:/ /:/ /\/__/::::| \:\/__/:/ \__\:\/__/:/ /:/_____|__|:|           \  \:\ /  /:/__/:/ /:/\:|
	//     \  \:\/\  \:\~~\~~\/\  \:\~~\__\/\  \:\/:/ /:/\  \:\~~\__\/\  \:\ /  /:/\  \:\/:::::/__/::::\            \  \:\  /:/\  \:\/:/~/:/
	//      \__\::/\  \:\  ~~~  \  \:\       \  \::/ /:/  \  \:\       \  \:\  /:/  \  \::/~~~~   ~\~~\:\            \  \:\/:/  \  \::/ /:/
	//      /__/:/  \  \:\       \  \:\       \  \:\/:/    \  \:\       \  \:\/:/    \  \:\         \  \:\            \  \::/    \  \:\/:/
	//      \__\/    \  \:\       \  \:\       \  \::/      \  \:\       \  \::/      \  \:\         \__\/             \__\/      \  \::/
	//                \__\/        \__\/        \__\/        \__\/        \__\/        \__\/                                       \__\/

	// `)

	db := db.New()
	h := cli.NewHandler(db)

	// Start with an implicit transaction ---- my opinion
	db.Begin()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		// if exit := handler.Process(line); exit {
		if h.Process(line) {
			break
		}
	}
}
