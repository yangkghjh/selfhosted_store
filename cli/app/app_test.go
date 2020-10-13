package app

import (
	"testing"

	"github.com/docker/cli/cli/compose/types"
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

var sampleAppMetadata = Metadata{
	Title:       "Yarr",
	Name:        "yarr",
	Description: "开源 RSS 阅读器，Go 实现，数据存储于 SQLite。",
	Categories:  []string{"Read"},
	Platform:    "linux",
	Note:        "通过 IP:7070 打开。",
}

func TestLoadApp(t *testing.T) {
	app, err := NewApp(sampleAppYaml)

	assert.NilError(t, err)
	assert.DeepEqual(t, app, &App{
		Type:     "container",
		Metadata: sampleAppMetadata,
	})
}

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

var sampleConfig = types.Config{
	Version: "3.3",
}

func TestLoadDockerCompose(t *testing.T) {
	app, err := NewApp(sampleAppYaml)
	assert.NilError(t, err)

	err = app.LoadDockerCompose(sampleDockerComposeYaml)
	assert.NilError(t, err)

	actual := app.Compose
	assert.Equal(t, actual.Version, sampleConfig.Version)
	assert.Equal(t, len(actual.Services), 1)
	assert.Equal(t, actual.Services[0].ContainerName, "yarr")
}
