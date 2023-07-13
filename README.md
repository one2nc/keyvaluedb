
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

## Story 1 (set, get, and delete commands)
Implement the following commands: SET, GET and DELETE
```
$ SET name foo
> OK

$ SET surname "foo bar"
> OK

$ GET name
> "foo"

$ DEL surname
> (integer) 1  

$ GET surname           
> (nil)

$ SET surname bar 
> OK
```

## Story 2 (incr and incrby commands)
Implement Basic Numeric Operations (INCR, INCRBY) with Error Handling
```
$ SET counter 0
> OK

$ INCR counter    
> (integer) 1

$ GET counter     
> "1"

$ INCRBY counter 10 
> (integer) 11

$ INCR foo          
> (integer) 1

$ INCRBY bar 21    
> (integer) 21
```

## Story 3 (multi, exec, and discard commands)
Implement the following commands: MULTI, EXEC, and DISCARD
### Case 1: Happy path
```
$ MULTI
> OK

$ INCR foo        
> QUEUED

$ SET bar 1 
> QUEUED

$ EXEC
> 1) (integer) 1
  2) OK
```
### Case 2: Discard
```
$ MULTI           
> OK

$ INCR foo        
> QUEUED

$ SET bar 1 
> QUEUED

$ DISCARD   
> OK

$ GET key1  
> (nil)
```
## Story 4 (compact command)
Implement a `COMPACT` command that outputs the current state of the data store. This is a custom command that the actual Redis server doesn’t implement.
###
### Example 1:
```
$ SET counter 10
> OK

$ INCR counter
> OK

$ INCR counter
> OK

$ SET foo bar
> OK

$ GET counter    
> "12"

$ INCR counter
> "13"

$ COMPACT        
> SET counter 13
  SET foo bar
```
### Example 2:
```
$ INCR counter      
> OK 

$ INCRBY counter 10 
> OK

$ GET counter       
> "11"

$ DEL counter      
> (integer) 1

$ COMPACT           
> (nil)
```
