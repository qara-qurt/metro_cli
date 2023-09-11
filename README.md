# METRO CLI

**METRO CLI** is a command-line tool for accessing the metro schedule in Almaty, Kazakhstan. It provides information on metro station schedules and allows users to view station-specific schedules and other options.

## Installation

To install METRO CLI, you need to have Go (Golang) installed. If you haven't already, you can download and install it from [the official Go website](https://golang.org/dl/).

Once Go is installed, you can install METRO CLI using the following command:

```bash
go get -u github.com/qara-qurt/metroCLI
```

Or you can copy this repository or download
```bash
cd ./metroCLI
./metroCLI.exe
```

## Options
METRO CLI provides the following command-line options:

--station value, -s value, --stat value: Specify the station for which you want to view the schedule (default: 0).
--all, -a: Show the full schedule table for the specified station (default: false).
--help, -h: Show help and usage information.

## Examples
Here are some examples of how to use METRO CLI:

View the schedule:
```bash
metroCLI
```

View the schedule for station number 5:
```bash
metroCLI --station 5
```
View the full schedule table for station number 2:
```bash
metroCLI --all --station 2
```
Display help and usage information:
```bash
metroCLI --help
```
