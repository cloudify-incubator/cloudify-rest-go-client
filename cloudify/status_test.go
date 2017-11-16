/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

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

package cloudify

import (
	tests "github.com/cloudify-incubator/cloudify-rest-go-client/cloudify/tests"
	"testing"
)

const statusResponce = `{
	"status": "running",
	"services": [{
		"instances": [],
		"display_name": "Cloudify Composer"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "LSB: Starts Logstash as a daemon.",
			"state": "running",
			"MainPID": 0,
			"Id": "logstash.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name":
		"Logstash"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "nginx - high performance web server",
			"state": "running",
			"MainPID": 1038,
			"Id": "nginx.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "Webserver"
	}, {
		"instances": [],
		"display_name": "Cloudify Stage"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "InfluxDB Service",
			"state": "running",
			"MainPID": 949,
			"Id": "cloudify-influxdb.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "InfluxDB"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "Cloudify AMQP InfluxDB Broker Service",
			"state": "running",
			"MainPID": 2680,
			"Id": "cloudify-amqpinflux.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "AMQP InfluxDB"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "RabbitMQ Service",
			"state": "running",
			"MainPID": 1697,
			"Id": "cloudify-rabbitmq.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "RabbitMQ"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "Cloudify Management Worker Service",
			"state": "running",
			"MainPID": 2683,
			"Id": "cloudify-mgmtworker.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "Celery Management"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "PostgreSQL 9.5 database server",
			"state": "running",
			"MainPID": 1029,
			"Id": "postgresql-9.5.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "PostgreSQL"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "Cloudify REST Service",
			"state": "running",
			"MainPID": 950,
			"Id": "cloudify-restservice.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "Manager Rest-Service"
	}, {
		"instances": [{
			"LoadState": "loaded",
			"Description": "Riemann Service",
			"state": "running",
			"MainPID": 2682,
			"Id": "cloudify-riemann.service",
			"ActiveState": "active",
			"SubState": "running"
		}],
		"display_name": "Riemann"
	}]
}`

// TestGetStatus - check GetStatus
func TestGetStatus(t *testing.T) {
	var conn tests.FakeClient
	conn.GetResponse = []byte(statusResponce)
	conn.GetError = nil
	cl := ClientFromConnection(&conn)
	status, err := cl.GetStatus()
	if err != nil {
		t.Error("Recheck error reporting")
	}
	if status.Status != "running" {
		t.Errorf("Recheck unmarshal for 'status' field '%s'", status.Status)
	}
	if status.Services[0].Status() != "unknown" {
		t.Errorf("Recheck unmarshal for 'status' in first service '%s'", status.Services[0].Status())
	}
	if status.Services[1].Status() != "running" {
		t.Errorf("Recheck unmarshal for 'status' in second service '%s'", status.Services[1].Status())
	}
}

const versionResponce = `{
	"date": null,
	"edition": "community",
	"version": "17.6.30",
	"build": null,
	"commit": null
}`

// TestGetVersion - check GetVersion
func TestGetVersion(t *testing.T) {
	var conn tests.FakeClient
	conn.GetResponse = []byte(versionResponce)
	conn.GetError = nil
	cl := ClientFromConnection(&conn)
	version, err := cl.GetVersion()
	if err != nil {
		t.Error("Recheck error reporting")
	}
	if version.Version != "17.6.30" {
		t.Errorf("Recheck unmarshal for 'version' field '%s'", version.Version)
	}
}
