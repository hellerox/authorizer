// +build integration

package main

import (
	cmd2 "authorizer/internal/root"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"authorizer/internal/app/service"
	"authorizer/internal/app/storage"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		name   string
		writer *bytes.Buffer
		db     service.Storage
	}{
		{"run",
			new(bytes.Buffer),
			&storage.InMemory{},
		},
		{"simple-run",
			new(bytes.Buffer),
			&storage.InMemory{},
		},
		{"double-creation",
			new(bytes.Buffer),
			&storage.InMemory{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := service.New(tt.db)

			input, err := os.Open("testdata/" + tt.name + ".in")
			if err != nil {
				assert.NoError(t, err)
			}

			cmd2.Execute(service, input, tt.writer)

			expected, err := ioutil.ReadFile("testdata/" + tt.name + ".out")
			if err != nil {
				assert.NoError(t, err)
			}

			assert.Equal(t, string(expected), tt.writer.String())
		})
	}
}
