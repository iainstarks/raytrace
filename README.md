# raytrace
Go code to try and implement the code in the prag prog Ray Tracing book

I've taken the feature file from the prag prog website
Created this code in my github, and created a go module based on it

Then installed godog
go get github.com/cucumber/godog/cmd/godog

This has added the godog command to ~/go/bin

Which I used to suggest some dummy test code by running
godog features/tuples.feature

It gives me the code to set up a feature context by setting up appropriate regexps to tie
my feature text to a load of functions
and creates dummy implementations for this functions - so what I then have to do is wire up
those dummy functions so that my scenarios pass.

Here's a snippet from the output:

Scenario: Reflecting a vector off a slanted surface # features/tuples.feature:146
Given v ← vector(0, -1, 0)
And n ← vector(√2/2, √2/2, 0)
When r ← reflect(v, n)
Then r = vector(1, 0, 0)

30 scenarios (30 undefined)
88 steps (88 undefined)
5.77941ms

You can implement step definitions for undefined steps with these snippets:

func aATuple(arg1, arg2, arg3, arg4, arg5, arg6 int) error {
return godog.ErrPending
}

func aIsAPoint() error {
return godog.ErrPending
}

This is a decent article about how to do godog:
https://semaphoreci.com/community/tutorials/how-to-use-godog-for-behavior-driven-development-in-go
