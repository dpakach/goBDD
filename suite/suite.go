package suite

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/parser"
	"github.com/dpakach/gorkin/token"
)

type StepDef struct {
	Token   token.Type
	Pattern string
	Action  reflect.Value
}

func (s *StepDef) Match(step object.Step) bool {
	if step.Token.Type != s.Token {
		return false
	}

	if step.StepText == s.Pattern {
		dataLen := len(step.Data)
		if step.Table != nil {
			dataLen += 1
		}
		if dataLen != s.Action.Type().NumIn() {
			return false
		}

		for i, val := range step.Data {
			if s.Action.Type().In(i) != reflect.TypeOf(val) {
				return false
			}
		}

		if step.Table != nil {
			if s.Action.Type().In(len(step.Data)) != reflect.TypeOf(step.Table) {
				return false
			}
		}
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
	if numArgs := s.Action.Type().NumIn(); numArgs != len(callArgs) {
		return errors.New(fmt.Sprintf("Number of arguments mismatched, expected %v, got %v", len(callArgs), numArgs))
	}

	s.Action.Call(callArgs)
	return nil
}

type Suite struct {
	steps []*StepDef
}

func NewSuite() *Suite {
	return &Suite{
		[]*StepDef{},
	}
}

func (s *Suite) Given(pattern string, action interface{}) {
	s.AddStep(token.GIVEN, pattern, action)
}

func (s *Suite) When(pattern string, action interface{}) {
	s.AddStep(token.WHEN, pattern, action)
}

func (s *Suite) Then(pattern string, action interface{}) {
	s.AddStep(token.THEN, pattern, action)
}

func (s *Suite) RunFeature(input string) {
	l := lexer.New(input)
	p := parser.New(l)

	s.Run(p.Parse())
}

func (s *Suite) Run(res *object.FeatureSet) {
	var background *object.Background
	fail := false
	for _, feat := range res.Features {
		if feat.Background != nil {
			background = feat.Background
		}

		for _, scenario := range feat.Scenarios {
			for _, sc := range scenario.GetScenarios() {
				err := s.RunScenario(sc, background)
				if err != nil {
					fail = true
					break
				}
			}
		}
	}

	if fail {
		os.Exit(1)
	}
}

func (c *Suite) AddStep(token token.Type, pattern string, action interface{}) error {
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

func (c *Suite) GetMatch(step object.Step) (*StepDef, error) {
	for _, stepDef := range c.steps {
		if stepDef.Match(step) {
			return stepDef, nil
		}
	}
	log.Printf("Could not find step definition for:\n %v %v\n%v\n", step.Token.Type.String(), step.StepText, getSnippet(step))
	return nil, errors.New("Could not find step definition")
}

func getSnippet(s object.Step) string {
	return fmt.Sprintf(`
	suite.%v.("%v", func(args...) {
		// Your code here
	})
	`, s.Token.Type, s.StepText)
}

func (c *Suite) RunScenario(scenario object.Scenario, background *object.Background) error {
	steps := []object.Step{}
	if background != nil && len(background.Steps) > 0 {
		for _, step := range background.Steps {
			steps = append(steps, step)
		}
	}
	if scenario.Steps != nil {
		for _, step := range scenario.Steps {
			steps = append(steps, step)
		}
	}

	fail := false
	for _, step := range steps {
		if _, err := c.GetMatch(step); err != nil {
			fail = true
		}
	}

	if fail {
		return errors.New("Undefined Steps")
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
