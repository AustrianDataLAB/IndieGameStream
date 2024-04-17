# AD: Angular as frontend framework
**Decider:** Kai Pietruska

## Context and Problem Statement
What framework do we use to implement our api?

## Decision Drivers
We want to use a framework to create the api client, to make our life easier.
It should be simple since most of us dont have further golang experience and the api won't have complex functionality.

## Considered Options
- Gin-Gonic
- Beego

## Decision Outcome
Gin, because it is lightweight.
MVC (beego) is too overkill for our use-case and we already use angular for our frontend.

## Consequences
tbd