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

package internal

var logFile = &File{
	Name: "log",
	imports: []string{
		`otlplogs "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/logs/v1"`,
	},
	testImports: []string{
		`"testing"`,
		``,
		`"github.com/stretchr/testify/assert"`,
		``,
		`otlplogs "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/logs/v1"`,
	},
	structs: []baseStruct{
		resourceLogsSlice,
		resourceLogs,
		instrumentationLibraryLogsSlice,
		instrumentationLibraryLogs,
		logSlice,
		logRecord,
	},
}

var resourceLogsSlice = &sliceStruct{
	structName: "ResourceLogsSlice",
	element:    resourceLogs,
}

var resourceLogs = &messageStruct{
	structName:     "ResourceLogs",
	description:    "// ResourceLogs is a collection of logs from a Resource.",
	originFullName: "otlplogs.ResourceLogs",
	fields: []baseField{
		resourceField,
		&sliceField{
			fieldMame:       "InstrumentationLibraryLogs",
			originFieldName: "InstrumentationLibraryLogs",
			returnSlice:     instrumentationLibraryLogsSlice,
		},
	},
}

var instrumentationLibraryLogsSlice = &sliceStruct{
	structName: "InstrumentationLibraryLogsSlice",
	element:    instrumentationLibraryLogs,
}

var instrumentationLibraryLogs = &messageStruct{
	structName:     "InstrumentationLibraryLogs",
	description:    "// InstrumentationLibraryLogs is a collection of logs from a LibraryInstrumentation.",
	originFullName: "otlplogs.InstrumentationLibraryLogs",
	fields: []baseField{
		instrumentationLibraryField,
		&sliceField{
			fieldMame:       "Logs",
			originFieldName: "Logs",
			returnSlice:     logSlice,
		},
	},
}

var logSlice = &sliceStruct{
	structName: "LogSlice",
	element:    logRecord,
}

var logRecord = &messageStruct{
	structName:     "LogRecord",
	description:    "// LogRecord are experimental implementation of OpenTelemetry Log Data Model.\n",
	originFullName: "otlplogs.LogRecord",
	fields: []baseField{
		&primitiveTypedField{
			fieldMame:       "Timestamp",
			originFieldName: "TimeUnixNano",
			returnType:      "TimestampUnixNano",
			rawType:         "uint64",
			defaultVal:      "TimestampUnixNano(0)",
			testVal:         "TimestampUnixNano(1234567890)",
		},
		traceIDField,
		spanIDField,
		&primitiveTypedField{
			fieldMame:       "Flags",
			originFieldName: "Flags",
			returnType:      "uint32",
			rawType:         "uint32",
			defaultVal:      `uint32(0)`,
			testVal:         `uint32(0x01)`,
		},
		&primitiveField{
			fieldMame:       "SeverityText",
			originFieldName: "SeverityText",
			returnType:      "string",
			defaultVal:      `""`,
			testVal:         `"INFO"`,
		},
		&primitiveTypedField{
			fieldMame:       "SeverityNumber",
			originFieldName: "SeverityNumber",
			returnType:      "SeverityNumber",
			rawType:         "otlplogs.SeverityNumber",
			defaultVal:      `SeverityNumberUNDEFINED`,
			testVal:         `SeverityNumberINFO`,
		},
		&primitiveField{
			fieldMame:       "Name",
			originFieldName: "Name",
			returnType:      "string",
			defaultVal:      `""`,
			testVal:         `"test_name"`,
		},
		bodyField,
		attributes,
		droppedAttributesCount,
	},
}

var bodyField = &messageField{
	fieldName:       "Body",
	originFieldName: "Body",
	returnMessage:   anyValue,
}
