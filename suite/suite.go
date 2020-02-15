package suite

import (
	"log"

	"github.com/dpakach/gorkin/lexer"
	"github.com/dpakach/gorkin/object"
	"github.com/dpakach/gorkin/parser"
	"github.com/dpakach/gorkin/token"
	"github.com/dpakach/goBDD/container"
)

type Suite struct {
	C *container.Container
}

func NewSuite() *Suite {
	return &Suite {
		container.NewContainer(),
	}
}

func (s *Suite)Given(pattern string, action interface{}) {
	s.C.AddStep(token.GIVEN, pattern, action)
}

func (s *Suite)When(pattern string, action interface{}) {
	s.C.AddStep(token.WHEN, pattern, action)
}

func (s *Suite)Then(pattern string, action interface{}) {
	s.C.AddStep(token.THEN, pattern, action)
}

func (s *Suite)RunFeature(input string) {
	l := lexer.New(input)
	p := parser.New(l)

	s.Run(p.Parse())
}

func (s *Suite)Run(res *object.FeatureSet) {
	var background *object.Background
	for _, feat := range res.Features {
		if feat.Background != nil {
			background = feat.Background
		}

		for _, scenario := range feat.Scenarios {
			for _, sc := range scenario.GetScenarios() {
				err := s.C.Run(sc, background)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
