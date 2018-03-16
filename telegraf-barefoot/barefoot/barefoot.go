// The MIT License (MIT)
//
// Copyright (c) 2017 AT&T
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package inputs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Metric struct {
	ChassisPort     int   `json:"chassis_port"`
	OctetsIn        int64 `json:"octets_in"`
	OctetsOut       int64 `json:"octets_out"`
	PacketsDropped  int64 `json:"packets_dropped_buffer_full"`
	PacketsIn       int64 `json:"packets_in"`
	PacketsOut      int64 `json:"packets_out"`
}

type Metrics []Metric

type Barefoot struct {
	Url                string
	lastTime           time.Time
	lastOctetsIn       [512]int64
	lastOctetsOut      [512]int64
	lastPacketsDropped [512]int64
	lastPacketsIn      [512]int64
	lastPacketsOut     [512]int64
}

var sampleConfig = `
  ## URL-prefix for Barefoot
  url = "http://localhost:8100/"
`

func (_ *Barefoot) Description() string {
	return `Gather Barefoot Metrics`
}

func (_ *Barefoot) SampleConfig() string {
	return sampleConfig
}

func (s *Barefoot) Gather(acc telegraf.Accumulator) error {
	defer func() {
		if r := recover(); r != nil {
			glog.Error("E! Problem reading from Barefoot")
			acc.AddFields("status", map[string]interface{}{"ready": false}, nil, time.Now())
		}
	}()

	now := time.Now()
	diffTime := now.Sub(s.lastTime).Seconds()

	var requestURL = fmt.Sprint(s.Url, "metrics")
	content, err := getContent(requestURL)
	if err != nil {
		glog.Error("Error talking to Barefoot:", err)
		acc.AddFields("status", map[string]interface{}{"ready": false}, nil, time.Now())
		return err
	}

	var metrics Metrics
	
	err = json.Unmarshal(content, &metrics)
	if err != nil {
		glog.Error("Error Umarshalling:", err)
		glog.Error("content:", content)
		return err
	}

	for _, metric := range metrics {
		tags := map[string]string{
			"port": fmt.Sprintf("port_%v", metric.ChassisPort),
		}

		var (
			metricsCount = 0
			octetsIn       int64 = 0
			octetsOut      int64 = 0
			packetsDropped int64 = 0
			packetsIn      int64 = 0
			packetsOut     int64 = 0
		)

		if metric.OctetsIn == 0 {
			metricsCount++
		} else if s.lastOctetsIn[metric.ChassisPort] != 0 {
			octetsIn = (metric.OctetsIn - s.lastOctetsIn[metric.ChassisPort]) / int64(diffTime)
			metricsCount++
		}
		s.lastOctetsIn[metric.ChassisPort] = metric.OctetsIn

		if metric.OctetsOut == 0 {
			metricsCount++
		} else if s.lastOctetsOut[metric.ChassisPort] != 0 {
			octetsOut = (metric.OctetsOut - s.lastOctetsOut[metric.ChassisPort]) / int64(diffTime)
			metricsCount++
		}
		s.lastOctetsOut[metric.ChassisPort] = metric.OctetsOut

		if metric.PacketsDropped == 0 {
			metricsCount++
		} else if s.lastPacketsDropped[metric.ChassisPort] != 0 {
			packetsDropped = (metric.PacketsDropped - s.lastPacketsDropped[metric.ChassisPort])
			metricsCount++
		}
		s.lastPacketsDropped[metric.ChassisPort] = metric.PacketsDropped

		if metric.PacketsIn == 0 {
			metricsCount++
		} else if s.lastPacketsIn[metric.ChassisPort] != 0 {
			packetsIn = (metric.PacketsIn - s.lastPacketsIn[metric.ChassisPort]) / int64(diffTime)
			metricsCount++
		}
		s.lastPacketsIn[metric.ChassisPort] = metric.PacketsIn

		if metric.PacketsOut == 0 {
			metricsCount++
		} else if s.lastPacketsOut[metric.ChassisPort] != 0 {
			packetsOut = (metric.PacketsOut - s.lastPacketsOut[metric.ChassisPort]) / int64(diffTime)
			metricsCount++
		}
		s.lastPacketsOut[metric.ChassisPort] = metric.PacketsOut

		if metricsCount == 5 {
			acc.AddGauge("ports", map[string]interface{}{"octets_in": octetsIn}, tags, now)
			acc.AddGauge("ports", map[string]interface{}{"octets_out": octetsOut}, tags, now)
			acc.AddGauge("ports", map[string]interface{}{"packets_dropped_buffer_full": packetsDropped}, tags, now)
			acc.AddGauge("ports", map[string]interface{}{"packets_in": packetsIn}, tags, now)
			acc.AddGauge("ports", map[string]interface{}{"packets_out": packetsOut}, tags, now)
		}
	}

	acc.AddFields("status", map[string]interface{}{"ready": true}, nil, now)

	s.lastTime = now

	return nil
}

func getContent(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func init() {
	inputs.Add("barefoot", func() telegraf.Input {
		return &Barefoot{}
	})
}
