![Alt text](https://img.shields.io/badge/version-production-green.svg)
# snap collector plugin - sessioninfo
Collects Paloalto firewall session info  

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [License](#license-and-authors)
4. [Releases](#Releases)
5. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building. See also [How to install Go language](http://ask.xmodulo.com/install-go-language-linux.html)

### Operating systems
Builds for: 
* Linux/amd64

### Installation
#### To build the plugin binary:
```
$ go get -u github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo
```
### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo/tree/master/examples).

## Documentation

### Collected Metrics

This plugin has the ability to gather the following metric:
[Available Metrics](METRICS.md)

### Example
Example running sessioninfo collector and writing data to an Influx database.

Load sessioninfo plugin
```
$ snaptel plugin load $GOPATH/bin/snap-plugin-collector-sessioninfo
```
List available plugins
```
$ snaptel plugin list
NAME                             VERSION         TYPE            SIGNED          STATUS          LOADED TIME
sessioninfo                      2               collector       false           loaded          Fri, 02 Dec 2016 15:00:51 EST
```
See available metrics for your system
```
$ snaptel metric list
```

Create a task manifest file and put firewall 'api' key and 'ip' address  (example files in [examples] (https://github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo/blob/master/examples/task.yml)). 
Do not change 'cmd' in task manifest.
```yaml
version: 1
schedule:
  type: "simple"
  interval: "30s"
max-failures: 15
workflow:
  collect:
    metrics:
      /pan/sessioninfo/Age_accel_tsf:           {}
      /pan/sessioninfo/Age_scan_ssf:            {}
      /pan/sessioninfo/Age_scan_thresh:         {}
      /pan/sessioninfo/Age_scan_tmo:            {}
      /pan/sessioninfo/Cps:                     {}
      /pan/sessioninfo/Dis_def:                 {}
      /pan/sessioninfo/Dis_tcp:                 {}
      /pan/sessioninfo/Dis_udp:                 {}
      /pan/sessioninfo/Icmp_unreachable_rate:   {}
      /pan/sessioninfo/Kbps:                    {}
      /pan/sessioninfo/Max_pending_mcast:       {}
      /pan/sessioninfo/Num_active:              {}
      /pan/sessioninfo/Num_bcast:               {}
      /pan/sessioninfo/Num_icmp:                {}
      /pan/sessioninfo/Num_installed:           {}
      /pan/sessioninfo/Num_max:                 {}
      /pan/sessioninfo/Num_mcast:               {}
      /pan/sessioninfo/Num_predict:             {}
      /pan/sessioninfo/Num_tcp:                 {}
      /pan/sessioninfo/Num_udp:                 {}
      /pan/sessioninfo/Pps:                     {}
      /pan/sessioninfo/Tmo_cp:                  {}
      /pan/sessioninfo/Tmo_def:                 {}
      /pan/sessioninfo/Tmo_icmp:                {}
      /pan/sessioninfo/Tmo_tcp:                 {}
      /pan/sessioninfo/Tmo_tcp_unverif_rst:     {}
      /pan/sessioninfo/Tmo_tcphalfclosed:       {}
      /pan/sessioninfo/Tmo_tcphandshake:        {}
      /pan/sessioninfo/Tmo_tcpinit:             {}
      /pan/sessioninfo/Tmo_tcptimewait:         {}
      /pan/sessioninfo/Tmo_udp:                 {}
      /pan/sessioninfo/Vardata_rate:            {}
      /pan/sessioninfo/Age_accel_thresh:        {}
    config:
      /pan/sessioninfo:
        api: ""
        ip: ""
        cmd: "&cmd=<show><session><info/></session></show>"
    tags:
      /pan/sessioninfo/:
        site: "DC1"
        id: "2"
    publish:
      -
        plugin_name: "influxdb"
        config:
          host: "localhost"
          port: 8086
          database: "test"
          user: "admin"
          password: "admin"
          https: false
          skip-verify: false
          retention: ten_weeks
```
Load sessioninfo plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-influxdb
```

Create a task:
```
$ snaptel task create -t task.yml
Using task manifest to create task
Task created
ID: 031c21b1-475b-41a6-8053-675fff2c9b9d
Name: Task-031c21b1-475b-41a6-8053-675fff2c9b9d
State: Running
```

List running tasks:
```
$ snaptel task list
ID                                       NAME                                            STATE           HIT     MISS    FAIL    CREATED                 LAST FAILURE
031c21b1-475b-41a6-8053-675fff2c9b9d     Task-031c21b1-475b-41a6-8053-675fff2c9b9d       Running         0       0       0       3:01PM 12-02-2016
```
Watch the task
```
$snaptel task watch 031c21b1-475b-41a6-8053-675fff2c9b9d
Watching Task (031c21b1-475b-41a6-8053-675fff2c9b9d):
NAMESPACE                                DATA            TIMESTAMP
/pan/sessioninfo/Age_accel_thresh        80              2017-03-30 09:05:41.416429245 -0400 EDT
/pan/sessioninfo/Age_accel_tsf           2               2017-03-30 09:05:41.416433457 -0400 EDT
/pan/sessioninfo/Age_scan_ssf            8               2017-03-30 09:05:41.416439235 -0400 EDT
/pan/sessioninfo/Age_scan_thresh         80              2017-03-30 09:05:41.416439847 -0400 EDT
/pan/sessioninfo/Age_scan_tmo            10              2017-03-30 09:05:41.416440558 -0400 EDT
/pan/sessioninfo/Cps                     1372            2017-03-30 09:05:41.416399708 -0400 EDT
/pan/sessioninfo/Dis_def                 60              2017-03-30 09:05:41.416412736 -0400 EDT
/pan/sessioninfo/Dis_tcp                 90              2017-03-30 09:05:41.416430182 -0400 EDT
/pan/sessioninfo/Dis_udp                 60              2017-03-30 09:05:41.416441334 -0400 EDT
/pan/sessioninfo/Icmp_unreachable_rate   200             2017-03-30 09:05:41.416430912 -0400 EDT
/pan/sessioninfo/Kbps                    38339           2017-03-30 09:05:41.416413635 -0400 EDT
/pan/sessioninfo/Max_pending_mcast       0               2017-03-30 09:05:41.416431774 -0400 EDT
/pan/sessioninfo/Num_active              94485           2017-03-30 09:05:41.41640127 -0400 EDT
/pan/sessioninfo/Num_bcast               1               2017-03-30 09:05:41.416414445 -0400 EDT
/pan/sessioninfo/Num_icmp                114             2017-03-30 09:05:41.416434235 -0400 EDT
/pan/sessioninfo/Num_installed           1671123632      2017-03-30 09:05:41.416415361 -0400 EDT
/pan/sessioninfo/Num_max                 4194302         2017-03-30 09:05:41.416418057 -0400 EDT
/pan/sessioninfo/Num_mcast               0               2017-03-30 09:05:41.416408843 -0400 EDT
/pan/sessioninfo/Num_predict             868             2017-03-30 09:05:41.416434954 -0400 EDT
/pan/sessioninfo/Num_tcp                 63706           2017-03-30 09:05:41.416432567 -0400 EDT
/pan/sessioninfo/Num_udp                 29886           2017-03-30 09:05:41.416442012 -0400 EDT
/pan/sessioninfo/Pps                     9480            2017-03-30 09:05:41.416409447 -0400 EDT
/pan/sessioninfo/Tmo_cp                  30              2017-03-30 09:05:41.416403356 -0400 EDT
/pan/sessioninfo/Tmo_def                 30              2017-03-30 09:05:41.416438484 -0400 EDT
/pan/sessioninfo/Tmo_icmp                6               2017-03-30 09:05:41.416420734 -0400 EDT
/pan/sessioninfo/Tmo_tcp                 3600            2017-03-30 09:05:41.416418952 -0400 EDT
/pan/sessioninfo/Tmo_tcp_unverif_rst     30              2017-03-30 09:05:41.416404675 -0400 EDT
/pan/sessioninfo/Tmo_tcphalfclosed       120             2017-03-30 09:05:41.41641025 -0400 EDT
/pan/sessioninfo/Tmo_tcphandshake        10              2017-03-30 09:05:41.416405749 -0400 EDT
/pan/sessioninfo/Tmo_tcpinit             5               2017-03-30 09:05:41.416407833 -0400 EDT
/pan/sessioninfo/Tmo_tcptimewait         15              2017-03-30 09:05:41.416419879 -0400 EDT
/pan/sessioninfo/Tmo_udp                 30              2017-03-30 09:05:41.416427242 -0400 EDT
/pan/sessioninfo/Vardata_rate            10485760        2017-03-30 09:05:41.416428421 -0400 EDT
```
Watch metrics in real-time using [Snap plugin for Grafana] (https://blog.raintank.io/using-grafana-with-intels-snap-for-ad-hoc-metric-exploration/) 
and use sessioninfo plugin for publishing ![Alt text](examples/grafana-sessioninfo.JPG "Metrics published to influxdb")

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE).

## Releases
### Version 2

Initial release

### TODO

## Acknowledgements
* Author: [@IrekRomaniuk](https://github.com/IrekRomaniuk/)

