/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package installer

import (
	"github.com/ca17/teamsedge/common"
	"github.com/ca17/teamsedge/config"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

var InstallScript = `#!/bin/bash -x
groupadd teamsedge
useradd teamsedge -g teamsedge -M -s /sbin/nologin
mkdir -p /var/teamsedge
chown -R teamsedge.teamsedge /var/teamsedge
chmod -R 700 /var/teamsedge
install -m 777 ./teamsedge /usr/local/bin/teamsedge 
chown teamsedge.teamsedge /etc/teamsedge.yaml 
test -d /usr/lib/systemd/system || mkdir -p /usr/lib/systemd/system
cat>/usr/lib/systemd/system/teamsedge.service<<EOF
[Unit]
Description=teamsedge
After=network.target

[Service]
Environment=GODEBUG=x509ignoreCN=0
LimitNOFILE=65535
LimitNPROC=65535
Username=teamsedge
ExecStart=/usr/local/bin/teamsedge

[Install]
WantedBy=multi-user.target
EOF

chmod 600 /usr/lib/systemd/system/teamsedge.service
systemctl enable teamsedge && systemctl daemon-reload

`

func InitConfig(config *config.AppConfig) error {
	// config.NBI.JwtSecret = common.UUID()
	cfgstr, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("/etc/teamsedge.yaml", cfgstr, 0644)
}

func Install(config *config.AppConfig) error {
	if !common.FileExists("/etc/teamsedge.yaml") {
		_ = InitConfig(config)
	}
	script := strings.ReplaceAll(InstallScript, "/var/teamsedge", config.System.Workdir)
	cmd := "/usr/local/bin/teamsedge"
	script = strings.ReplaceAll(InstallScript, "/usr/local/bin/teamsedge", cmd)
	_ = ioutil.WriteFile("/tmp/teamsedge_install.sh", []byte(script), 0777)

	// 创建用户&组
	if err := exec.Command("/bin/bash", "/tmp/teamsedge_install.sh").Run(); err != nil {
		return err
	}
	return os.Remove("/tmp/teamsedge_install.sh")
}

func Uninstall() {
	_ = os.Remove("/etc/teamsedge.yaml")
	_ = os.Remove("/usr/lib/systemd/system/teamsedge.service")
	_ = os.Remove("/usr/local/bin/teamsedge")
}
