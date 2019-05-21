# `ccfetch`

A tool for fetching code coverage and other statistics from team city.
Expectes a file `secrets.json` in the current directory with the format

```json
{
  "host": "http://your.teamcity.server/",
  "username": "snazz-meister",
  "password": "topsecret"
}
```
