# goBDD
A BDD based test runner for golang

## Getting Started

1.  First install goBDD in your project or test infrastructure
```bash
go get -u github.com/dpakach/goBDD
```
2. Create a feature file for your project. For writing feature files we use the [gherkin language](https://cucumber.io/docs/gherkin/). Use [this](https://cucumber.io/docs/gherkin/reference/) reference from cucumber.io to write your feature files.
Here is a simple example of a feature file
```gherkin
Feature: Some feature

    Scenario: An example scenario
        When I do "BDD"
        Then my project has less bugs
```
The feature should match the feature you want to test and the scenario should be an example of that feature described using the gherkin specification.

3. Create a go program that contains the definitions for your gherkin steps. You can write your step definitions to do any action or assertions for your tests.
Here is an example of a simple step definitions file.
``` go
package main

import (
    "fmt"

    "github.com/dpakach/goBDD/runner"
    "github.com/dpakach/goBDD/suite"
)

func main() {
    s := suite.NewSuite()

    s.When("I do {{s}}", func(task string) {
        fmt.Printf("I am doing %v\n", task)
    })

    s.Then("my project has less bugs", func() {
        fmt.Println("congrats! Your project has less bugs now.")
    })
    runner.Run(s)
}

```

Here we are just printing some text on the stdOut but you can do anything you want in the step definitions. You can call functions from your project and test their results, make api calls and test the api response or even integrate a selenium driver and perform UI tests on your project.

4. Now to run the tests go to the terminal and run the go program.
``` bash
go run main.go [path/to/your/feature/files]
```

## License

Copyright (c) 2020 Dipak Acharya

Licensed under [MIT License](https://github.com/dpakach/goBDD/blob/master/LICENSE)
