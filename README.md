# Filebeat Multi-line Tester

## Description

A small Go application to test a multi-line regex to be used with the filebeat `multiline` config option

## Building application

Run `make build` in the root of this application repository.

## Usage

- `-f` : File containing multi-line string to test (string, default: empty)
- `-n` : Negate the pattern matching (bool, default: `true`)
- `-p` : Multi-line regex pattern to use for the matching (string, default: empty)
- `-y` : Specify a filebeat prospector yaml config, which overrides the `-f`, `-n`, and `-p` flags (default: "")
- `-t` : The pattern matching backend to use for the regex, either `1` for `regexp.CompilePOSIX`, `2` for `regexp.Compile`, or `3` for `match.Matcher`.  This will depend on the filebeat version you're using. (int, default: 3)
- `-s` : Print the individual lines (default: `false`)
- `-v` : Prints current version and exits

### Patern matching backends

Depending on the filebeat version, different backends were used to match the regular
expression.  Here are the version details:
  - version 1.x						= regexp.CompilePOSIX
  - version 5.0 to 5.3		= regexp.Compile
  - version 5.3 and up 		= match.Matcher


## Example

```
./filebeat-multiline-tester -p "^=[A-Z]+|^$" -f teststring.txt
```
or
```
./filebeat-multiline-tester -y sample_configs/conf-diskusage.yml -t 2
```

## Credits

This code base is a adaptation of the code sample which Elastic provides as a testing mechanism within the Go Playground.

https://play.golang.org/p/uAd5XHxscu


## License

Coverted under the [MIT license](LICENSE.md). 