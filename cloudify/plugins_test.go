/*
Copyright (c) 2018 GigaSpaces Technologies Ltd. All rights reserved

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

const pluginsResponce = `{
	"items": [{
		"distribution_release": "core",
		"supported_py_versions": ["py27"],
		"uploaded_at": "2018-03-02T18:31:51.065Z",
		"archive_name": "cloudify_utilities_plugin-1.5.0-py27-none-linux_x86_64-centos-Core.wgn",
		"package_version": "1.5.0",
		"package_name": "cloudify-utilities-plugin",
		"distribution_version": "7.3.1611",
		"tenant_name": "default_tenant",
		"excluded_wheels": [],
		"created_by": "admin",
		"distribution": "centos",
		"package_source": "../cloudify-utilities-plugin/",
		"private_resource": false,
		"file_server_path": "",
		"resource_availability": "tenant",
		"visibility": "tenant",
		"supported_platform": "linux_x86_64",
		"wheels": [
			"paramiko-2.4.0-py2.py3-none-any.whl",
			"Jinja2-2.10-py2.py3-none-any.whl",
			"xmltodict-0.11.0-py2.py3-none-any.whl",
			"requests-2.18.4-py2.py3-none-any.whl",
			"requests_toolbelt-0.8.0-py2.py3-none-any.whl",
			"pyasn1-0.4.2-py2.py3-none-any.whl",
			"certifi-2018.1.18-py2.py3-none-any.whl",
			"chardet-3.0.4-py2.py3-none-any.whl",
			"idna-2.6-py2.py3-none-any.whl",
			"urllib3-1.22-py2.py3-none-any.whl",
			"six-1.11.0-py2.py3-none-any.whl",
			"asn1crypto-0.24.0-py2.py3-none-any.whl",
			"enum34-1.1.6-py2-none-any.whl",
			"cloudify_plugins_common-3.4.2-py2-none-any.whl",
			"cloudify_rest_client-4.0-py2-none-any.whl",
			"cloudify_utilities_plugin-1.5.0-py2-none-any.whl",
			"pycrypto-2.6.1-cp27-cp27mu-linux_x86_64.whl",
			"PyYAML-3.12-cp27-cp27mu-linux_x86_64.whl",
			"pika-0.9.14-py2-none-any.whl",
			"networkx-1.8.1-py2-none-any.whl",
			"proxy_tools-0.1.0-py2-none-any.whl",
			"bottle-0.12.7-py2-none-any.whl",
			"bcrypt-3.1.4-cp27-cp27mu-linux_x86_64.whl",
			"cryptography-2.1.4-cp27-cp27mu-linux_x86_64.whl",
			"PyNaCl-1.2.1-cp27-cp27mu-linux_x86_64.whl",
			"MarkupSafe-1.0-cp27-cp27mu-linux_x86_64.whl",
			"cffi-1.11.4-cp27-cp27mu-linux_x86_64.whl",
			"ipaddress-1.0.19-py2-none-any.whl",
			"pycparser-2.18-py2.py3-none-any.whl"
		],
		"id": "54ecf152-1e01-41be-9eda-e849f78e6eea",
		"yaml_url_path": "plugin:cloudify-utilities-plugin?version=1.5.0&distribution=centos"
	}],
	"metadata": {
		"pagination": {
			"total": 1,
			"offset": 0,
			"size": 100
		}
	}
}`

// TestGetPlugins - check GetPlugins
func TestGetPlugins(t *testing.T) {
	var conn tests.FakeClient
	conn.GetResponse = []byte(pluginsResponce)
	conn.GetError = nil
	cl := ClientFromConnection(&conn)
	plugins, err := cl.GetPlugins(map[string]string{"id": "54ecf152-1e01-41be-9eda-e849f78e6eea"})
	if err != nil {
		t.Error("Recheck error reporting")
	}
	tests.AssertEqual(t, plugins.Items[0].ID, "54ecf152-1e01-41be-9eda-e849f78e6eea",
		"Recheck unmarshal for 'id' field in plugin '%s'", plugins.Items[0].ID)
}
