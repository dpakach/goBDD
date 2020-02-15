package main

import (
	"fmt"

	"github.com/dpakach/goBDD/suite"
	"github.com/dpakach/gorkin/object"

	"github.com/dpakach/goBDD/runner"
)

func main() {

	s := suite.NewSuite()
	s.When("I run background", func() {
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println("Running Background")
	})

	s.Then("I am happy", func(table object.Table) {
		//fmt.Println("I am happy")
		//fmt.Println(table)
	})

	s.When("I do something", func(table object.Table) {
		//fmt.Println("I do something")
		//fmt.Println(table)
	})

	s.Then("something happens", func() {
		//fmt.Println("I do something")
	})

	s.When("i do something {{s}}", func(task string) {
		fmt.Println(task)
		//fmt.Printf("I am doing %v task\n", task)
	})

	s.Then("something {{s}} happens", func(res string) {
		//fmt.Printf("%v is happening\n", res)
	})
	runner.Run(s)
}