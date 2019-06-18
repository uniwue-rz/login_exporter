# Login Exporter

This exporter uses the xpath, simple page downloader and regex matching
to check login to different type of web apps. Then exports the data
for prometheus to scrape. As logging in an expensive operation it uses
cache to serve the last cached data to Prometheus. The checks on 
different web pages is done with the help of go routines.


## Running

To run this application simply call it in the terminal

```bash
./login_exporter
```

## Configuration


This application uses the chromedriver on the given machine, to use
that you need to install google chrome and the driver on the machine
running the application.

### On Ubuntu

### On MacOS

## Development

This application is open source and can be extended.

### Build

To build this application you need to install the needed requirements
first:

```bash
go get
go build -o ./login_exporter
```