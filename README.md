# In-Memory Database

Language used: Go (Golang)

## How to Run

1. Build the program (from [cmd/main.go](https://github.com/taytaybear/Go-Database/blob/master/cmd/main.go)):
   ```bash
   go build -o database cmd/main.go
   ```
2. Run the program:
   ```bash
   ./database
   ```
3. Enter commands one per line. To exit, type `END` or use <kbd>Ctrl</kbd>+<kbd>D</kbd> (EOF).

## Assumptions
- [x] All input commands are well-formed and space-delimited.
- [x] Commands and keys are case-sensitive
- [x] All values are treated as strings.
- [x] Data is stored in memory only and is lost when the program exits.
- [x] No external dependencies are used; only the Go standard library.
