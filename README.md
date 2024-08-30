# Go Smartschool API wrapper
A simple Go package for the Smartschool API.

## Yes. But.. Why?
I always wanted to do stuff with Smartschool, like pairing it with Home Assistant so it can read aloud my new messages.

## How to use

Simply run
```shell
go add github.com/alessiodam/gosmartschool
```
and start using it!


## Examples

for this, I'll use examples/auth_check

1. Make a .env file with the following:
```dotenv
DOMAIN=school.smartschool.be
PHPSESSID=your phpsessid cookies' value
PID=your pid cookies' value
```
2. Run `go run ./examples/auth_check/`
3. You should get a message "Authenticated!" if your .env is correct.
