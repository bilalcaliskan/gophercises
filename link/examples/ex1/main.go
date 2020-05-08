package main

import (
	"fmt"
	"gophercises/link"
	"strings"
)

var exampleHtml = `
<html>
<body>
<h1>Hello!</h1>
<a href="/other-page">A link to another page</a>
<a href="/page-two">A link to a second <span> some span </span> page</a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}