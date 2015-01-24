[![Build Status](https://travis-ci.org/tiago4orion/DataGen.svg?branch=master)](https://travis-ci.org/tiago4orion/DataGen) [![Coverage Status](https://coveralls.io/repos/tiago4orion/DataGen/badge.svg)](https://coveralls.io/r/tiago4orion/DataGen)

# Data Generator

This tool tries to easy the process of generate random datasets for database infrastructure's tests.

# Installation

```bash
go get github.com/tiago4orion/DataGen
```

# Usage

You need a YAML file defining your "schema" or "document". See below:

```YAML
format: csv
filename: company.csv
fields:
  - id:
      type: "string"
      min: 14
      max: 14
      chars: "0-9"
  - company_name:
      type: "string"
      min: 3
      max: 256
      chars: " A-Za-z"
  - address:
      type: "string"
      min: 3
      max: 256
      chars: " A-Za-z"
```
This file only specifies the data types that will be generated.
Then, run the code below:

```bash
$GOPATH/bin/DataGen
-c is required...
-n is required...
Usage: ./DataGen [options]:
           --help, -h: Displays this help
    --data-config, -c: Data configuration file
    --output-file, -o: Output file
         --format, -f: Output format
 --number-records, -n: Number of records
     --concurrent, -t: Number of concurrent routines
     
$GOPATH/bin/DataGen -c sample.yaml -n 10000 -t 8 -o data 
Workers status: 100.00%
```
See the generated files:
```bash
wc -l file-out_*
   1246 file-out_0.csv
   1254 file-out_1.csv
   1253 file-out_2.csv
   1256 file-out_3.csv
   1241 file-out_4.csv
   1255 file-out_5.csv
   1249 file-out_6.csv
   1246 file-out_7.csv
  10000 total
  ```
  
  Very simple tool, but can be useful ;)
  

