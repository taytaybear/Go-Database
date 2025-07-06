// main handler for the database application
// accepts input from the user and processes it
// returns true if the user wants to exit the application
// returns false if the user wants to continue
package cli

import (
	"fmt"
	"os"
	"strings"

	"go-database/internal/db"
)

type Handler struct {
	db *db.DB
}

func NewHandler(d *db.DB) *Handler {
	return &Handler{db: d}
}

func (h *Handler) Process(input string) bool {
	input = strings.TrimSpace(input)
	args := strings.Fields(input)
	if len(args) == 0 {
		return false
	}

	cmd := args[0]

	switch cmd {
	case "SET":
		if len(args) != 3 {
			return false
		}
		h.db.Set(args[1], args[2])
	case "GET":
		if len(args) != 2 {
			return false
		}
		if v, ok := h.db.Get(args[1]); ok {
			fmt.Println(v)
		} else {
			fmt.Println("NULL")
		}
		os.Stdout.Sync()
	case "UNSET":
		if len(args) != 2 {
			return false
		}
		h.db.Unset(args[1])
	case "NUMEQUALTO":
		if len(args) != 2 {
			return false
		}
		fmt.Println(h.db.NumEqualTo(args[1]))
		os.Stdout.Sync()
	case "BEGIN":
		if len(args) != 1 {
			return false
		}
		h.db.Begin()
	case "ROLLBACK":
		if len(args) != 1 {
			return false
		}
		if !h.db.Rollback() {
			fmt.Println("NO TRANSACTION")
			os.Stdout.Sync()
		}
	case "COMMIT":
		if len(args) != 1 {
			return false
		}
		if !h.db.Commit() {
			fmt.Println("NO TRANSACTION")
			os.Stdout.Sync()
		}
	case "END":
		if len(args) != 1 {
			return false
		}
		return true
	default:
		return false
	}
	return false
}
