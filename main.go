package main

import "github.com/article/app"

func main() {
	a := app.App{}
	a.Initialize()
	a.Run()
}
