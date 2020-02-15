# 4chan-dl
![loc](https://sloc.xyz/github/nektro/4chan-dl)
[![license](https://img.shields.io/github/license/nektro/4chan-dl.svg)](https://github.com/nektro/4chan-dl/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)
[![paypal](https://img.shields.io/badge/donate-paypal-009cdf)](https://paypal.me/nektro)
[![circleci](https://circleci.com/gh/nektro/4chan-dl.svg?style=svg)](https://circleci.com/gh/nektro/4chan-dl)
[![release](https://img.shields.io/github/v/release/nektro/4chan-dl)](https://github.com/nektro/4chan-dl/releases/latest)
[![goreportcard](https://goreportcard.com/badge/github.com/nektro/4chan-dl)](https://goreportcard.com/report/github.com/nektro/4chan-dl)
[![codefactor](https://www.codefactor.io/repository/github/nektro/4chan-dl/badge)](https://www.codefactor.io/repository/github/nektro/4chan-dl)

Media downloader for 4chan.org.

## Prerequisites
- Golang 1.12+

## Installing
```sh
$ go get -v -u github.com/nektro/4chan-dl
```

## Usage
```
Usage of ./4chan-dl:
  -b, --board stringArray   /--board/ to download.
      --concurrency int     Maximum number of simultaneous downloads. (default 10)
      --save-dir string     Path to a directory to save to.
```
Example:
```
$ 4chan-dl --board wg --save-dir ./downloads/
```

## Built With
- https://github.com/nektro/go-util
- https://github.com/spf13/pflag
- https://github.com/valyala/fastjson

## License
MIT
