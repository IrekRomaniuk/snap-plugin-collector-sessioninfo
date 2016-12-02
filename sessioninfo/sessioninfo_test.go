package sessioninfo

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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
	"fmt"
)

func TestSessioninfoPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName )
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

func TestSessioninfoCollector_CollectMetricsollectMetrics(t *testing.T) {
	cfg := setupCfg("paste api key here","10.34.2.21","&cmd=<show><session><info/></session></show>")

	Convey("Sessioninfo collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 1)

		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"pan", "sessioninfo", "num-active"),
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
				fmt.Println(m.Namespace()[2].Value,m.Data())
				So(m.Namespace()[2].Value, ShouldBeIn, "num-active")
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