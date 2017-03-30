package sessioninfo

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

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	"encoding/xml"

)

const (
	//put your api key and ip address here
	api = "XXX"
	ip  = "10.34.2.21"
)

func TestSessioninfoPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName)
		So(meta.Version, ShouldResemble, pluginVersion)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})
	Convey("Create Sessioninfo Collector", t, func() {
		collector := New()
		Convey("So Sessioninfo collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So Sessioninfo collector should be of Sessioninfo type", func() {
			So(collector, ShouldHaveSameTypeAs, &SessioninfoCollector{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /pan/sessioninfo", func() {
				conf := configPolicy.Get([]string{"pan", "sessioninfo"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 3)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "api", "ip", "cmd")
					switch rule.Name {
					case "api":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "ip":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "cmd":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					}
				}
			})
		})
	})
}
// mock your 'get_page()' function here
func mock_get_page(url string) ([]byte, error){
	return []byte(sessioninfo_response), nil
}
// test session info fetch by looking at Kbps field
func TestSessioninfoFetch(t *testing.T) {
	cmd := "&cmd=<show><session><info/></session></show>"
	d := NewDownloader(mock_get_page)
	//d := NewDownloader(get_page)
	htmlData, _ := d.download("https://" + ip + "/esp/restapi.esp?type=op" + cmd + "&key=" + api)
	var session Session
	xml.Unmarshal(htmlData, &session)

	//fmt.Println(session.Result.Kbps, get_sub_field(session,"Result","Kbps"))
	Convey("Sessioninfo download Kbps", t, func() {
		So(get_sub_field(session,"Result","Kbps"), ShouldEqual, 1519)
	})
}

func TestSessioninfoCollector(t *testing.T) {
	cfg := setupCfg(api, ip, "&cmd=<show><session><info/></session></show>")
	//fmt.Println(api, ip)
	Convey("Sessioninfo collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 33)

		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				{
					Namespace_: core.NewNamespace(
						"pan", "sessioninfo", "Num_active"),
					Config_: cfg.ConfigDataNode,
				},
			}
			//fmt.Println(mts[0].Config().Table())
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "pan")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "sessioninfo")
			for _, m := range metrics {
				//fmt.Println("\n",m.Namespace()[2].Value,m.Data_)
				So(m.Namespace()[2].Value, ShouldEqual, "Num_active")
				t.Log(m.Namespace()[2].Value, m.Data())
			}
		})
	})
}

func setupCfg(api string, ip string, cmd string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("api", ctypes.ConfigValueStr{Value: api})
	node.AddItem("ip", ctypes.ConfigValueStr{Value: ip})
	node.AddItem("cmd", ctypes.ConfigValueStr{Value: cmd})
	return plugin.ConfigType{ConfigDataNode: node}
}

const sessioninfo_response = `<response status="success"><result>
<tmo-udp>30</tmo-udp>
<tcp-nonsyn-rej>False</tcp-nonsyn-rej>
<tmo-tcp>3600</tmo-tcp>
<pps>1124</pps>
<num-max>4194302</num-max>
<age-scan-thresh>80</age-scan-thresh>
<tmo-tcphalfclosed>120</tmo-tcphalfclosed>
<num-active>102448</num-active>
<dis-def>60</dis-def>
<num-mcast>0</num-mcast>
<icmp-unreachable-rate>200</icmp-unreachable-rate>
<tmo-tcptimewait>15</tmo-tcptimewait>
<age-scan-ssf>8</age-scan-ssf>
<vardata-rate>10485760</vardata-rate>
<age-scan-tmo>10</age-scan-tmo>
<tmo-tcpinit>5</tmo-tcpinit>
<dp>*.dp0</dp>
<dis-tcp>90</dis-tcp>
<num-udp>29570</num-udp>
<tmo-icmp>6</tmo-icmp>
<max-pending-mcast>0</max-pending-mcast>
<age-accel-thresh>80</age-accel-thresh>
<tmo-tcphandshake>10</tmo-tcphandshake>
<oor-action>drop</oor-action>
<tmo-def>30</tmo-def>
<age-accel-en>True</age-accel-en>
<age-accel-tsf>2</age-accel-tsf>
<hw-offload>True</hw-offload>
<num-icmp>114</num-icmp>
<num-predict>0</num-predict>
<tmo-cp>30</tmo-cp>
<strict-checksum>True</strict-checksum>
<tmo-tcp-unverif-rst>30</tmo-tcp-unverif-rst>
<num-bcast>0</num-bcast>
<ipv6-fw>True</ipv6-fw>
<num-installed>4142667796</num-installed>
<num-tcp>72139</num-tcp>
<dis-udp>60</dis-udp>
<cps>750</cps>
<kbps>1519</kbps>
</result></response>`
