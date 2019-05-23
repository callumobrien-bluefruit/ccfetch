# `ccfetch`

A tool for fetching code coverage and other statistics from team city.
Fetches the specified properties from the statistics for the latest
build of the specified build type and prints them in JSON to `stdout`.

## Usage

```
$ ccfetch -host <host> -id <id> -props <props>
```
where
- `<host>` is the address of your TeamCity server (default
  "http://127.0.0.1/")
- `<id>` is the ID of the build type to fetch properties for
- `<props>` is a colon-seperated list of the names of the properties to
  fetch

Expectes a file `secrets.json` in the current directory with the format
```
{
  "username": "snazz-meister",
  "password": "topsecret"
}
```
