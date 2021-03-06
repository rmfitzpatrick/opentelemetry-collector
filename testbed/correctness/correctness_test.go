// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package correctness

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/service/defaultcomponents"
	"go.opentelemetry.io/collector/testbed/testbed"
)

var correctnessResults testbed.TestResultsSummary = &testbed.CorrectnessResults{}

func TestMain(m *testing.M) {
	testbed.DoTestMain(m, correctnessResults)
}

func TestTracingGoldenData(t *testing.T) {
	tests, err := LoadPictOutputPipelineDefs("testdata/generated_pict_pairs_traces_pipeline.txt")
	assert.NoError(t, err)
	processors := map[string]string{
		"batch": `
  batch:
    send_batch_size: 1024
`,
	}
	for _, test := range tests {
		test.TestName = fmt.Sprintf("%s-%s", test.Receiver, test.Exporter)
		test.DataSender = ConstructTraceSender(t, test.Receiver)
		test.DataReceiver = ConstructReceiver(t, test.Exporter)
		t.Run(test.TestName, func(t *testing.T) {
			testWithTracingGoldenDataset(t, test.DataSender, test.DataReceiver, test.ResourceSpec, processors)
		})
	}
}

func testWithTracingGoldenDataset(
	t *testing.T,
	sender testbed.DataSender,
	receiver testbed.DataReceiver,
	resourceSpec testbed.ResourceSpec,
	processors map[string]string,
) {
	dataProvider := testbed.NewGoldenDataProvider(
		"../../internal/goldendataset/testdata/generated_pict_pairs_traces.txt",
		"../../internal/goldendataset/testdata/generated_pict_pairs_spans.txt",
		161803)
	factories, err := defaultcomponents.Components()
	assert.NoError(t, err)
	runner := testbed.NewInProcessCollector(factories, sender.GetCollectorPort())
	validator := testbed.NewCorrectTestValidator(dataProvider)
	config := CreateConfigYaml(sender, receiver, processors, "traces")
	configCleanup, cfgErr := runner.PrepareConfig(config)
	assert.NoError(t, cfgErr)
	defer configCleanup()
	tc := testbed.NewTestCase(
		t,
		dataProvider,
		sender,
		receiver,
		runner,
		validator,
		correctnessResults,
	)
	defer tc.Stop()

	tc.SetResourceLimits(resourceSpec)
	tc.EnableRecording()
	tc.StartBackend()
	tc.StartAgent("--metrics-level=NONE")

	tc.StartLoad(testbed.LoadOptions{
		DataItemsPerSecond: 1024,
		ItemsPerBatch:      1,
	})

	duration := time.Second
	tc.Sleep(duration)

	tc.StopLoad()

	tc.WaitForN(func() bool { return tc.LoadGenerator.DataItemsSent() == tc.MockBackend.DataItemsReceived() },
		duration, "all data items received")

	tc.StopAgent()

	tc.ValidateData()
}
