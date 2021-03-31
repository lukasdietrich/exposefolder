# exposefolder

A oneline web server for development purposes.

![Screenshot](https://user-images.githubusercontent.com/1913775/113221224-e6c5e780-9284-11eb-977f-e3c95720ac1d.png)

## Motivation

When I need to quickly serve a folder via http for whatever reason I usually use
`python3 -m http.server`.

Apart from developing a web application, which is a one way street from the server to the browser,
I sometimes wished you could easily upload files without any hassle.

`exposefolder` serves a folder with file listings, when there is no index file present. On a listing
you can just drag and drop file into the browser and they appear, just like magic ✨.

## Usage

```sh
$ exposefolder -h
Usage of exposefolder:
  -foldername string
	Path of the folder to expose (default ".")
  -port int
	Port to bind the web server on (default 8080)
```

```sh
$ exposefolder -foldername myproject -port 1234
2021/04/01 00:48:58 serving "myproject" on ":1234"
```

## Installation

Using Go:

```sh
$ go install github.com/lukasdietrich/exposefolder@master
```
