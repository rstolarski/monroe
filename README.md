# Marilyn Diptych

This app generates procedurally short animation of Marilyn Monroe faces being printed.
This was a task at Lodz Universoty of Technology in one of classes.

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

* Marta Osowicz - for proposing first version of this idea
* Krzysztof Guzek - for guidance during class
