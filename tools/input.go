package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskYesNo(s string) (bool, error) {
	var keep string
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(s + "([y]/n)")
		v, _ := reader.ReadString('\n')
		keep = strings.Trim(v, "\n")
		if keep == "" || keep == "y" || keep == "n" {
			break
		} else {
			fmt.Printf("Invalid input: %s\n", keep)
		}
	}

	return keep != "n", nil
}
