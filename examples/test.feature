@coolTag
Feature: test
	Background:
		When I run background
		Then I am happy
			| also | with |
			| a	| table |

	@tag
	Scenario: example scenario
		When I do something
			| also | with |
			| a	| table |
		Then something happens

	Scenario Outline: another example scenario
		When i do something "<task>"
		Then something "good" happens
			| a	| table |
		Examples:
			| task |
			| good |
			| bad  |
			| okay |
