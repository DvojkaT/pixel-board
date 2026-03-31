package name

import "math/rand"

type Name string

const (
	CAPYBARA Name = "капибара"
	SNAKE    Name = "змея"
	MELON    Name = "арбуз"
	BANANA   Name = "банан"
	BOAR     Name = "кабан"
	KANGAROO Name = "кенгуру"
)

var Names = []Name{CAPYBARA, SNAKE, MELON, BANANA, BOAR, KANGAROO}

func Random() Name {
	return Names[rand.Intn(len(Names))]
}
