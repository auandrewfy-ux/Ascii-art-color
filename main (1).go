package main

import (
	"fmt"
	"os"
	"strings"
)

/*
Here we define the colors our program understands.

Each color name maps to an ANSI escape code.
ANSI codes are special sequences that tell the terminal
to change how text is displayed (for example: color, bold, etc.).

Example:
"\033[31m" → switch to red text
"\033[0m"  → reset terminal formatting back to normal
now var is a keyword in Go used to declare a variable. Meaning "I want to create a variable in memory" and the name of our variable of course is 'color'.
our assignment operator '=' which is simply saying; take the value on the right and store it in the variable on the left.
map[string]string describes the structure of the map.

A map in Go stores data in key value pairs (like a dictionary).

The first "string" inside the brackets [string] represents the
type of the key, meaning the keys must be text.

The second "string" represents the type of the value, meaning
the values must also be text.

So map[string]string means:
a map where text keys point to text values.

Example:
"red" -> "\033[31m"

Here "red" is the key (the color name) and "\033[31m"
is the value (the ANSI escape code that tells the terminal
to display text in red).

ASCII, Escaping, and ANSI codes:
ASCII is a system that assigns a number to every character.
- Example: 'A' = 65, 'space' = 32, ESC = 27.
Escape sequences:
- Some characters are special or non-printable.
- We use a backslash \ to "escape" them in a string.
- Example: "\n" = newline, "\t" = tab, "\033" = ESC (ASCII 27).

Why we're escaping "\033":
- "\033" represents the ESC (escape) character, which cannot be typed directly.
- ESC tells the terminal that a **command follows**, not normal text.

Why we need it in my program:
- ANSI color codes use ESC + command, e.g., "\033[31m" = red text.
- The terminal interprets ESC[31m as "start red color".
- "\033[0m" resets the terminal back to normal formatting.
								      
So in short:
- We escape because the ESC character is non-printable.
- ESC signals the terminal that a formatting/color command follows.
- This allows my/the program to change text colors safely and then reset.

*/

var colors = map[string]string{
	"black":  "\033[30m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"white":  "\033[37m",
	"reset":  "\033[0m",
}

/*
This map stores the available ASCII-art banners.

The key is the banner name that the user can type,
and the value is the actual file that contains the ASCII patterns.

Example command:
go run . "hello" shadow

This means we must load the file "shadow.txt".
*/
var banners = map[string]string{
	"standard":   "standard.txt",
	"shadow":     "shadow.txt",
	"thinkertoy": "thinkertoy.txt",
}

/*
If the user runs the program incorrectly,
we print this usage message to guide them.
*/
func usage() {
	fmt.Println(`Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"`)
}

/*
This function reads a banner file from disk.

Example banner files:
standard.txt
shadow.txt
thinkertoy.txt

Each banner file contains the ASCII representation
of every printable character (from ASCII 32 to 126).

The file is split into lines so we can easily
extract the 8 rows that represent each character.

func readBanner(name string) ([]string, error)

This defines a function called readBanner that reads an ASCII banner from a file.

- `func` declares a function.
- `readBanner` the functions name (what we call to use it).
- `(name string)` input parameter: the name of the banner to read (e.g., "shadow"). Must be text.
- `([]string, error)` output: a slice of strings (each line of the banner) and an error (if something goes wrong).

Example call:

lines, err := readBanner("standard")

- `lines` will hold the ASCII banner lines.
- `err` will hold an error if the file does not exist or cannot be read.

*/
func readBanner(name string) ([]string, error) {

	// Check if the banner name exists in our banner list
	file, ok := banners[name]
	if !ok {
		return nil, fmt.Errorf("invalid banner")
	}

	// Read the entire banner file
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

// Split file contents into individual lines

return strings.Split(string(data), "\n"), nil
}

// string(data) simply says; take the raw bytes from the file and turn them into normal, readable text.
// Imagine you open a wrapped letter and unfold it to see the words.
// strings.Split(..., "\n") Cut that text into separate lines wherever there is a newline character.
// Think of slicing a long scroll into strips so each line can be read individually.
// return ... , nil; Give these strips (lines) back to whoever called the function, and also tell them there is no error.
// The messenger comes back to you and says: “Here is the banner, all ready, and nothing went wrong.”

// Extra comment by Mr curious: 
// The nil here only happens if everything went right.
// if something goes wrong while reading the file, this line does not even run.
// Why? Because before this line, Mr Curious has already checked for errors.

// so the difference between this line 👇🏼 and my my last readfile line is ;
/* data, err := os.ReadFile(file)
if err != nil {
    return nil, err
	}

	// err != nil; “something went wrong”
	return nil, err; stop the function and send the error back immediately
	The function never reaches the strings.Split line if reading failed.
*/

/*
This function extracts the ASCII-art pattern
for a single character.

Each character in the banner file occupies 8 lines.

The index calculation works like this:

ASCII characters start at code 32 (space).
So we subtract 32 to find the correct position.

Each character block takes 9 lines in the file
(8 lines of drawing + 1 empty separator line).
*/
func asciiChar(r rune, banner []string) []string {

	index := (int(r) - 32) * 9

	// Return the 8 rows representing this character
	return banner[index+1 : index+9]
}

/*
This function determines which characters
in the text should be colored.

It returns a boolean slice where:

true  → this character should be colored
false → leave it normal
*/
func findPositions(text, sub string) []bool {

	pos := make([]bool, len(text))

	/*
		If no substring was provided,
		the instructions say we must color
		the entire string.
	*/
	if sub == "" {
		for i := range pos {
			pos[i] = true
		}
		return pos
	}

	/*
		Here we search the text for occurrences
		of the substring.

		Whenever we find a match, we mark all
		the letters of that match as true.
	*/
	for i := 0; i <= len(text)-len(sub); i++ {
		if text[i:i+len(sub)] == sub {
			for j := 0; j < len(sub); j++ {
				pos[i+j] = true
			}
		}
	}

	return pos
}

/*
This function prints the ASCII art.

It combines:
- ASCII conversion
- substring detection
- color application
*/
func printAscii(text, substring, colorCode, bannerName string) {

	// Load the banner file
	banner, err := readBanner(bannerName)
	if err != nil {
		fmt.Println("Error reading banner")
		return
	}

	/*
		The program must support newline characters.

		So we split the input text by "\n"
		and print each line separately.
	*/
	lines := strings.Split(text, "\n")

	for _, line := range lines {

		// Determine which characters should be colored
		pos := findPositions(line, substring)

		/*
			Each ASCII character has 8 rows.

			So we print row 1 of every letter,
			then row 2 of every letter, etc.
		*/
		for row := 0; row < 8; row++ {

			for i, r := range line {

				// Ignore non printable ASCII
				if r < 32 || r > 126 {
					continue
				}

				// Get the ASCII drawing for this character
				char := asciiChar(r, banner)

				// If this position should be colored
				if colorCode != "" && pos[i] {
					fmt.Print(colorCode + char[row] + colors["reset"])
				} else {
					fmt.Print(char[row])
				}
			}

			// Move to next row of ASCII output
			fmt.Println()
		}
	}
}

/*
This function is responsible for interpreting
the command line arguments.

It determines:
- whether a color flag was provided
- whether a substring was provided
- which banner should be used
*/
func parse(args []string) (color, substring, text, banner string, ok bool) {

	// Default banner if none is specified
	banner = "standard"

	if len(args) == 0 {
		return
	}

	/*
		Check if the last argument is a banner name.
		If so, we use that banner and remove it
		from the argument list.
	*/
	last := args[len(args)-1]

	if _, exists := banners[last]; exists {
		banner = last
		args = args[:len(args)-1]
	}

	/*
		If there is only one argument left,
		it must be the text to print.
	*/
	if len(args) == 1 {
		text = args[0]
		ok = true
		return
	}

	/*
		Otherwise the first argument must
		be the color option.
	*/
	option := args[0]

	if !strings.HasPrefix(option, "--color=") {
		return
	}

	// Extract the color name from the flag
	colorName := strings.TrimPrefix(option, "--color=")

	colorCode, exists := colors[colorName]

	if !exists {
		return
	}

	color = colorCode

	/*
		If there are only two arguments,
		then the user wants to color the
		entire string.
	*/
	if len(args) == 2 {
		text = args[1]
		ok = true
		return
	}

	/*
		If there are three arguments,
		then we have a substring and a string.
	*/
	if len(args) == 3 {
		substring = args[1]
		text = args[2]
		ok = true
		return
	}

	return
}

/*
Main function:
The entry point of the program.

It parses arguments and then prints
the ASCII art result.
*/
func main() {

	color, substring, text, banner, ok := parse(os.Args[1:])

	if !ok {
		usage()
		return
	}

	printAscii(text, substring, color, banner)
}
