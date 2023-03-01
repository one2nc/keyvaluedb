
Write a basic key-value DB implementation (think redis 0.1 version). Program will accept DB commands as inputs and process them by creating DB structure in memory.

# Expectations:
---
- Choose any programming language.
- Test drive the code (use TDD). If not possible, write tests after the code is written. Code with zero tests will not be reviewed :-)
- Write Readme describing the Assumptions / Technical decisions etc.
- Aim of this exercise is to write modular and extensible code. It's okay if it is NOT highly performant, as we don't plan to use it in production to replace Redis.
- As with all things in software, the requirements will change over time. Good modular code is open for extension and closed for modifications (Open-Closed pricinple in SOLID). Aim to write such code.
- Don't get fancy, keep things simple and stupid. When in doubt, make reasonable assumptions and document them in the Readme.

Sample input is of format: `COMMAND ARGS...`

## Story 1: SET key-values and then GET them. You can also DEL a key.
```
SET key1 value1
SET key2 value2
GET key1 // should return value1
DEL key2 // deletes key2 and its value
GET key2 // should return nil
SET key2 newvalue2
```

## Story 2: Basic numeric operations and error cases
```
SET counter 0
INCR counter // increments a "counter" key, if present and returns incremented value, in this case: 1
GET counter // returns 1
INCRBY counter 10 // increment by 10, returns 11
INCR foo // automatically creates a new key called "foo" and increments it by 1 and thus returns 1
```

## Story 3: Command spanning multiple sub-commands
### Case 1: Happy path
```
MULTI // starts a multi line commands
INCR foo // queues this command, doesn't execute it immediately
SET key1 value1 // queues this command, doesn't execute it immediately
EXEC // execute all queued commands and returns output of all commands in an array, thus returns: [1 value1]
```
### Case 2: Discard
```
MULTI // starts a multi line commands
INCR foo // queues this command, doesn't execute it immediately
SET key1 value1 // queues this command, doesn't execute it immediately
DISCARD // discard all queued commands
GET key1 // returns nil as key1 doesn't exists
```
## Story 4: Generate compacted command output (think materialized view of current DB state)
### Example 1:
```
SET counter 10
INCR counter
INCR counter
SET foo bar
GET counter // returns 12
INCR counter
COMPACT // this should return following output
SET counter 13
SET foo bar
```
### Example 2:
```
INCR counter // returns 1
INCRBY counter 10 // returns 11
GET counter // returns 11
DEL counter // deletes counter
COMPACT // this should compact to empty output as there's no keys present in the DB
