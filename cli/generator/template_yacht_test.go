package generator

import (
	"testing"

	"github.com/yankghjh/store/cli/app"
	"gotest.tools/assert"
)

var sampleAppYaml = []byte(`
type: container
title: Yarr
name: yarr
description: 开源 RSS 阅读器，Go 实现，数据存储于 SQLite。
categories:
  - Read
platform: linux
note: "通过 IP:7070 打开。"
`)

var sampleDockerComposeYaml = []byte(`
version: '3.3'
services:
  yarr:
    container_name: yarr
    image: yangkghjh/yarr:latest
    ports:
      - 7070:7070/tcp # WebUI
      - 7071:8080 # AdminUI
      - 7072:8081/udp
    restart: always
    volumes:
      - "!data/yarr:/data"
    environment:
      - "SHARE=yacht;/mount"
`)

var sampleAppMetadata = app.Metadata{
	Title:       "Yarr",
	Name:        "yarr",
	Description: "开源 RSS 阅读器，Go 实现，数据存储于 SQLite。",
	Categories:  []string{"Read"},
	Platform:    "linux",
	Note:        "通过 IP:7070 打开。",
}

func Test_app2YachtTemplate(t *testing.T) {
	app, err := app.NewApp(sampleAppYaml)
	assert.NilError(t, err)
	err = app.LoadDockerCompose(sampleDockerComposeYaml)
	assert.NilError(t, err)

	actual, err := app2YachtTemplate(app)
	assert.NilError(t, err)

	assert.DeepEqual(t, actual, &YachtTemplate{
		Type:     1,
		Metadata: sampleAppMetadata,

		Image:         "yangkghjh/yarr:latest",
		RestartPolicy: "unless-stopped",
		Ports: []map[string]string{
			map[string]string{
				"7070": "7070:7070/tcp",
				"8080": "7071:8080/tcp",
				"8081": "7072:8081/udp",
			},
		},
		Volumes: []YachtVolumeConfig{
			YachtVolumeConfig{
				Container: "/data",
				Bind:      "!data/yarr",
			},
		},
		Environment: []YachtEnvironmentConfig{
			YachtEnvironmentConfig{
				Name:    "SHARE",
				Label:   "SHARE",
				Default: "yacht;/mount",
			},
		},
	})
}
