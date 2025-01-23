package main

import "studentgit.kata.academy/ponomarenko.100299/go-petstore/run"

func main() {
	app := run.NewApp()
	app.Bootstrap().Run()
}
