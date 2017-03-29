/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sessioninfo

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"github.com/PuerkitoBio/goquery"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"encoding/xml"
	"reflect"
)

const (
	vendor        = "pan"
	fs            = "sessioninfo"
	pluginName    = "sessioninfo"
	pluginVersion = 2
	pluginType    = plugin.CollectorPluginType
)

/*func init() {
}*/

type PageGetter func(url string) ([]byte, error)
type Downloader struct {
	get_page PageGetter
}
func NewDownloader(pg PageGetter) *Downloader {
	return &Downloader{get_page: pg}
}
type Session struct {
	XMLName xml.Name `xml:"response"`
	Status  string `xml:"status,attr""`
	Result Result `xml:"result"`
}
var (
	metricNames = []string{
		"Tmo_udp","Tmo_tcp","Pps","Num_max","Age_scan_thresh","Tmo_tcphalfclosed","Num_active","Dis_def",
		"Num_mcast","Icmp_unreachable_rate","Tmo_tcptimewait", "Age_scan_ssf","Vardata_rate","Age_scan_tmo",
		"Tmo_tcpinit","Dis_tcp","Num_udp","Tmo_icmp","Max_pending_mcast","Age_accel_thresh","Tmo_tcphandshake",
		"Tmo_def","Age_accel_tsf","Num_icmp","Num_predict","Tmo_cp","Tmo_tcp_unverif_rst","Num_bcast",
		"Num_installed", "Num_tcp","Dis_udp","Cps","Kbps",
	}
)
type Result struct {
	Tmo_udp string `xml:"tmo-udp"`
	Tmo_tcp string `xml:"tmo-tcp"`
 	Pps string `xml:"pps"`
	Num_max string `xml:"num-max"`
	Age_scan_thresh string `xml:"age-scan-thresh"`
	Tmo_tcphalfclosed string `xml:"tmo-tcphalfclosed"`
	Num_active string `xml:"num-active"`
	Dis_def string `xml:"dis-def"`
	Num_mcast string `xml:"num-mcast"`
	Icmp_unreachable_rate string `xml:"icmp-unreachable-rate"`
	Tmo_tcptimewait string `xml:"tmo-tcptimewait"`
	Age_scan_ssf string `xml:"age-scan-ssf"`
	Vardata_rate string `xml:"vardata-rate"`
	Age_scan_tmo string `xml:"age-scan-tmo"`
	Tmo_tcpinit string `xml:"tmo-tcpinit"`
	Dis_tcp string `xml:"dis-tcp"`
	Num_udp string `xml:"num-udp"`
	Tmo_icmp string `xml:"tmo-icmp"`
	Max_pending_mcast string `xml:"max-pending-mcast"`
	Age_accel_thresh string `xml:"age-accel-thresh"`
	Tmo_tcphandshake string `xml:"tmo-tcphandshake"`
	Tmo_def string `xml:"tmo-def"`
	Age_accel_tsf string `xml:"age-accel-tsf"`
	Num_icmp string `xml:"num-icmp"`
	Num_predict string `xml:"num-predict"`
	Tmo__cp string `xml:"tmo-cp"`
	Tmo_tcp_unverif_rst string `xml:"tmo-tcp-unverif-rst"`
	Num_bcast string`xml:"num-bcast"`
	Num_installed string `xml:"num-installed"`
	Num_tcp string `xml:"num-tcp"`
	Dis_udp string `xml:"dis-udp"`
	Cps string `xml:"cps"`
	Kbps string `xml:"kbps"`
}

type SessioninfoCollector struct {
}

func New() *SessioninfoCollector {
sessioninfo := &SessioninfoCollector{}
return sessioninfo
}

/*  CollectMetrics collects metrics for testing.

CollectMetrics() will be called by Snap when a task that collects one of the metrics returned from this plugins
GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

The output is the collected metrics as plugin.Metric and an error.
*/
func (sessioninfo *SessioninfoCollector) CollectMetrics(mts []plugin.MetricType) (metrics []plugin.MetricType, err error) {
//var err error
var (
	api string
	ip  string
	cmd string
	session Session
)
conf := mts[0].Config().Table()
apiConf, ok := conf["api"]
if !ok || apiConf.(ctypes.ConfigValueStr).Value == "" {
	return nil, fmt.Errorf("api key missing from config, %v", conf)
} else {
	api = apiConf.(ctypes.ConfigValueStr).Value
}
ipConf, ok := conf["ip"]
if !ok || ipConf.(ctypes.ConfigValueStr).Value == "" {
	return nil, fmt.Errorf("ip address missing from config, %v", conf)
} else {
	ip = ipConf.(ctypes.ConfigValueStr).Value
}
cmdConf, ok := conf["cmd"]
if !ok || cmdConf.(ctypes.ConfigValueStr).Value == "" {
	return nil, fmt.Errorf("cmd missing from config, %v", conf)
} else {
	cmd = cmdConf.(ctypes.ConfigValueStr).Value
}

d := NewDownloader(get_page)
htmlData, err := d.download("https://" + ip + "/esp/restapi.esp?type=op" + cmd + "&key=" + api)
if err != nil {
	return nil, fmt.Errorf("Error collecting metrics: %v", err)
}

xml.Unmarshal(htmlData, &session)
	//https://github.com/intelsdi-x/snap-plugin-lib-go/blob/master/examples/snap-plugin-collector-rand/rand/rand.go
for _, mt := range mts {//idx
	ns := mt.Namespace()
	//val := session.Result.Num_active //tmp until marshalling finished
	//val := get_sub_field(session,"Result","Num_active")
	val := get_sub_field(session,"Result",mt.Namespace()[2].Value)
	//fmt.Println("DEBUG: ",idx, val,ns,mt.Namespace()[2].Value)
	if err != nil {
		return nil, fmt.Errorf("Error collecting metrics: %v", err)
	}
	//fmt.Println(val)
	metric := plugin.MetricType{
		Namespace_: ns,
		Data_:      val,
		Timestamp_: time.Now(),
	}
	metrics = append(metrics, metric)
}
return metrics, nil
}
//fetching url and return html, can be mocked for testing
func get_page(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	resp.Body.Close()
	return htmlData, nil
}
//dependency injection for fetching url func get_page(url)
func (d *Downloader) download(url string) ([]byte,error) {
	content, err := d.get_page(url)
	if err != nil {
		return []byte{},err
	}
	return content, nil
}
//Get struct field value
func get_field(v Session, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}
//Get named substruct field value
func get_sub_field(v Session, field1 string, field2 string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field1).FieldByName(field2)
	return f.String()
}

/*
GetMetricTypes returns metric types for testing.
GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
available metrics).

Config info is passed in. This config information would come from global config snap settings.

The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (sessioninfo *SessioninfoCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
mts := []plugin.MetricType{}
for _, metricName := range metricNames {
	mts = append(mts, plugin.MetricType{
		Namespace_: core.NewNamespace("pan", "sessioninfo", metricName),
	})
}
return mts, nil
}

// GetConfigPolicy returns plugin configuration
func (sessioninfo *SessioninfoCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
c := cpolicy.New()
rule0, _ := cpolicy.NewStringRule("api", true)
rule1, _ := cpolicy.NewStringRule("ip", true)
rule2, _ := cpolicy.NewStringRule("cmd", true)
cp := cpolicy.NewPolicyNode()
cp.Add(rule0)
cp.Add(rule1)
cp.Add(rule2)
c.Add([]string{"pan", "sessioninfo"}, cp)
return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
return plugin.NewPluginMeta(
	pluginName,
	pluginVersion,
	pluginType,
	[]string{plugin.SnapGOBContentType}, //[]string{},
	[]string{plugin.SnapGOBContentType},
	plugin.Unsecure(true),
	plugin.ConcurrencyCount(1),
)
}
