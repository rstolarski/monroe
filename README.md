# Project Title

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.



## Installation

### Prerequisites

#### [FFmpeg](https://www.ffmpeg.org/) - Library used to output final video file

1. Download ffmpeg from their website
2. Export ffmpeg.exe from bin directory into directory known to environment variables

#### [Golang](https://golang.org/) - Programming language used to write this application

1.  Install golang from their website

### Clone

- Clone this repo to your local machine using `https://github.com/rtropisz/monroe`

### Setup
#### Windows
Build application using Golang.
```Shell
$ go build -o monroe.exe -v
```
Run script to output file.
```shell
$ .\generateFrames.bat
```

#### Unix
Build application using Golang.
```Shell
$ go build -o monroe -v
```
Run script to output file.
```shell
$ bash generateFrames.sh
```
---

## Authors

* **Ryszard Tropisz**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Krzysztof Guzek
