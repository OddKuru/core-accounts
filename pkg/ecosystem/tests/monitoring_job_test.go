package grpc_job

import (
	"net/http"
	"testing"

	ayaka "github.com/OddKuru/core-accounts/pkg/core"
	"github.com/OddKuru/core-accounts/pkg/ecosystem"
	"github.com/stretchr/testify/assert"
)

func TestMonitoringJobBuilder(t *testing.T) {
	t.Run("Should correctly build monitoring job", func(t *testing.T) {
		address := "localhost:1000"
		job, err := ecosystem.NewMonitoringJobBuilder[Container]().
			Address(address).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, address, job.Address())
	})

	t.Run("Should correctly failed building monitoring job", func(t *testing.T) {
		job, err := ecosystem.NewMonitoringJobBuilder[Container]().
			Build()

		assert.Error(t, err)
		assert.Nil(t, job)
	})

	t.Run("Should correctly custom mux", func(t *testing.T) {
		address := "localhost:1000"
		job, err := ecosystem.NewMonitoringJobBuilder[Container]().
			Address(address).
			Mux(http.NewServeMux()).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, address, job.Address())
	})
}

func TestMonitoringJobSignature(t *testing.T) {
	address := "localhost:10101"

	job, err := ecosystem.NewMonitoringJobBuilder[Container]().
		Address(address).
		Build()

	assert.NoError(t, err)

	ayaka.NewApp(&ayaka.Options[Container]{
		Name:        "aya",
		Description: "kekw",
		Version:     "0.0.1",
		Container:   Container{},
	}).WithJob(ayaka.JobEntry[Container]{
		Key: "xd",
		Job: job,
	})
}
