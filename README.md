# snap-plugin-collector-sessioninfo
collect Paloalto firewall session info

### setup
$go get -u github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo
$snaptel plugin load $GOPATH/bin/snap-plugin-collector-sessioninfo
$cp $GOPATH/src/github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo/examples/task.yml .
$snaptel task create -t task.yml


#### testing
cd $GOPATH//src/github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo/
go test
cd $GOPATH//src/github.com/IrekRomaniuk/snap-plugin-collector-sessioninfo/sessioninfo
go test
