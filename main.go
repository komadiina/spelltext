package main

import (
	"fmt"

	"github.com/komadiina/spelltext/utils/colors"
)

func main() {
	fmt.Println(colors.Paint("Hello world!", colors.Green))
}
