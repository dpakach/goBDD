package container

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/token"
)

type StepDef struct {
	Token token.Type
	Pattern string
	Action reflect.Value
}

func (s *StepDef) Match(step object.Step) bool {
	if step.Token.Type != s.Token {
		return false
	}

	if step.StepText == s.Pattern {
		return true
	}
	return false
}

func (s *StepDef) Run(args ...interface{}) error {
	callArgs := []reflect.Value{}
	for _, arg := range args {
		a := reflect.ValueOf(arg)
		callArgs = append(callArgs, a)
	}
	s.Action.Call(callArgs)
	return nil
}

type Container struct {
	steps []*StepDef
}

func NewContainer() *Container {
	return &Container{
		[]*StepDef{},
	}
}

func (c *Container) AddStep(token token.Type, pattern string, action interface{}) error {
	v := reflect.ValueOf(action)
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected handler to be func, but got: %T", action))
	}
	for _, stepDef := range c.steps {
		if stepDef.Pattern == pattern {
			return errors.New("Step Definition already exists")
		}
	}
	c.steps = append(c.steps, &StepDef{token, pattern, v})
	return nil
}

func (c *Container) GetMatch(step object.Step) (*StepDef, error) {
	for _, stepDef := range c.steps {
		if stepDef.Match(step) {
			return stepDef, nil
		}
	}
	log.Printf("Could not find step definition for: %v %v", step.Token, step.StepText)
	return nil, errors.New("Could not find step definition")
}

func (c *Container) Run(scenario object.Scenario, background *object.Background) error {
	steps := []object.Step{}
	if background.Steps != nil {
		for _, step := range background.Steps {
			steps = append(steps, step)
		}
	}
	if scenario.Steps != nil {
		for _, step := range scenario.Steps {
			steps = append(steps, step)
		}
	}
	for _, step := range steps {
		stepDef, err := c.GetMatch(step)
		if err != nil {
			return err
		}
		args := []interface{}{}
		for _, arg := range step.Data {
			i, err := strconv.Atoi(arg)
			if err == nil {
				args = append(args, i)
			} else {
				args = append(args, arg)
			}
		}
		if step.Table != nil && len(step.Table) > 0 {
			args = append(args, step.Table)
		}
		stepDef.Run(args...)
	}

	return nil
}

