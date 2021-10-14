

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("Enter text: \n")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(input)
}