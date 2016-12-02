# snap collector plugin - session info
This plugin collects Paloalto firewalls session info  

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [License](#license-and-authors)
4. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### To build the plugin binary:
```
$ go get -u github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo
```
### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/IrekRomaniuk/snap-plugin-collector-ping/blob/master/README.md#examples).

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/raintank/snap-plugin-collector-ping/blob/master/METRICS.md).

### Example
Example running ping collector and writing data to a file.

Load ping plugin
```
$ $SNAP_PATH/bin/snaptel plugin load snap-plugin-collector-sessioninfo
```
See available metrics for your system
```
$ $SNAP_PATH/bin/snaptel metric list
```

Create a task manifest file  (exemplary files in [examples/tasks/] (https://github.com/raintank/snap-plugin-collector-ping/blob/master/examples/)):
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/raintank/ping/*": {}
            },
            "config": {
            	"/raintank/ping": {
            		"hostname": "127.0.0.1"
            	}
            },
            "process": null,
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_ping"
                    }
                }
            ]
        }
    }
}
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snaptel plugin load $SNAP_PATH/plugin/snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create a task:
```
$ $SNAP_PATH/bin/snaptel task create -t examples/tasks/task.yml
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop previously created task:
```
$ $SNAP_PATH/bin/snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@IrekRomaniuk](https://github.com/IrekRomaniuk/)

