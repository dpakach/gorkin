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
		Then something happens

	Scenario Outline: another example scenario
		When i do something <task>
		Then something "good" happens

		Examples:
			| task |
			| good |
			| bad  |
		
