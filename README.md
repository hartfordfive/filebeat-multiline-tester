# Filebeat Multi-line Tester

## Description

A small Go application to test a multi-line regex to be used with the filebeat `multiline` config option

## Building application

Run `make all` in the root of this application repository.

## Usage

`-f` : File containing multi-line string to test (default: "")
`-n` : Negate the pattern matching (default: true)
`-p` : Multi-line regex pattern to use for the matching (default: "")
`-v` : Prints current version and exits

## Example

./multiline-test -p "^=[A-Z]+|^$" -f teststring.txt

## Credits

This code base is a adaptation of the code sample which Elastic provides a as a testing mechanism within the Go Playground.

https://play.golang.org/p/uAd5XHxscu

## License

Coverted under the [MIT license](LICENSE.md). 