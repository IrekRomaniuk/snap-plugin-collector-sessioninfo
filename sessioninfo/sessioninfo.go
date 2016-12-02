/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"strings"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"time"
)

const (
	vendor        = "pan"
	fs            = "sessioninfo"
	pluginName    = "sessioninfo"
	pluginVersion = 1
	pluginType    = plugin.CollectorPluginType
)

/*func init() {
}*/
var (
	metricNames = []string{
		"num-active",
	}
)


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
		ip string
		cmd string
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
        //fmt.Println("https://" + ip + "/esp/restapi.esp?type=op" + cmd + "&key=" + api)
	htmlData, err := getHTML("https://" + ip + "/esp/restapi.esp?type=op" + cmd + "&key=" + api)
	if err != nil { return nil, fmt.Errorf("Error collecting metrics: %v", err) }
	//fmt.Println(htmlData)
	for _, mt := range mts {
		ns := mt.Namespace()

		val, err := parseSessionInfo("num-active", htmlData)
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

func getHTML (url string ) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil { return "", err }
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil { return "", err }
	resp.Body.Close()
	return string(htmlData), nil
}
//HTML parse should go to snap-plugin-processor ?
func parseSessionInfo (tag string, htmlData string) (string, error) {
	htmlCode := strings.NewReader(htmlData)
	doc, err := goquery.NewDocumentFromReader(htmlCode)
	if err != nil { return "", err }
	s := doc.Find(tag).Text()
	return s, nil
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
		[]string{plugin.SnapGOBContentType},//[]string{},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(1),
	)
}
