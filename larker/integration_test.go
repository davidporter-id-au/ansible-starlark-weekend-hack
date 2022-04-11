package larker

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYamlConversion(t *testing.T) {
	expected := `# Generated code. Do not edit
clusters:
  - cluster-1:
      persistence:
        numHistoryShards: 8192
        defaultStore: caas-default
        visibilityStore: caas-visibility
      stats:
        exportInterval: 500ms
        exporter:
          m3:
            env: production
            hostPort: 127.0.0.1:9052
            service: cluster-1
      dynamicconfig:
        client: dynamic-configurator
        dynamic-configurator:
          namespaces: cluster-1
        applicationidentifier: application-server
        cachedir: /var/cache/dynamic-configurator-config
        iswatchfileenabled: "true"
  - cluster-2:
      persistence:
        numHistoryShards: 16384
        defaultStore: caas-default
        visibilityStore: caas-visibility
      stats:
        exportInterval: 500ms
        exporter:
          m3:
            env: production
            hostPort: 127.0.0.1:9052
            service: cluster-2
      dynamicconfig:
        client: dynamic-configurator
        dynamic-configurator:
          namespaces: cluster-2
        applicationidentifier: application-server
        cachedir: /var/cache/dynamic-configurator-config
        iswatchfileenabled: "true"
`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	l := New()

	res, err := l.Main(ctx, "testdata/vars.star")
	if err != nil {
		log.Fatal(err)
	}
	assert.NoError(t, err)

	assert.Equal(t, expected, res.YAMLConfig)
}
