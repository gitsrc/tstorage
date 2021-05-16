package tstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_metric_insertPoint(t *testing.T) {
	tests := []struct {
		name   string
		metric metric
		point  *DataPoint
		want   []DataPoint
	}{
		{
			name: "the first insertion",
			metric: metric{
				name:      "metric-a",
				points:    []DataPoint{},
				lastIndex: -1,
			},
			point: &DataPoint{
				Timestamp: 1,
				Value:     0,
			},
			want: []DataPoint{
				{
					Timestamp: 1,
					Value:     0,
				},
			},
		},
		{
			name: "insert in the middle",
			metric: metric{
				name: "metric-a",
				points: []DataPoint{
					{
						Timestamp: 1,
					},
					{
						Timestamp: 3,
					},
				},
				lastIndex: 1,
			},
			point: &DataPoint{
				Timestamp: 2,
			},
			want: []DataPoint{
				{
					Timestamp: 1,
				},
				{
					Timestamp: 2,
				},
				{
					Timestamp: 3,
				},
			},
		},
		{
			name: "insert into the last",
			metric: metric{
				name: "metric-a",
				points: []DataPoint{
					{
						Timestamp: 1,
					},
					{
						Timestamp: 2,
					},
				},
				lastIndex: 1,
			},
			point: &DataPoint{
				Timestamp: 3,
			},
			want: []DataPoint{
				{
					Timestamp: 1,
				},
				{
					Timestamp: 2,
				},
				{
					Timestamp: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.metric.insertPoint(tt.point)
			assert.Equal(t, tt.want, tt.metric.points)
		})
	}
}

func TestSelectAll(t *testing.T) {
	tests := []struct {
		name            string
		memoryPartition memoryPartition
		want            []Row
	}{
		{
			name: "single data point for single metric",
			memoryPartition: func() memoryPartition {
				m := memoryPartition{}
				m.metrics.Store("metric1", &metric{
					name: "metric1",
					points: []DataPoint{
						{
							Timestamp: 1,
							Value:     0.1,
						},
					},
					lastIndex: 0,
				})
				return m
			}(),
			want: []Row{
				{
					//MetricName: "metric1",
					DataPoint: DataPoint{
						Timestamp: 1,
						Value:     0.1,
					},
				},
			},
		},
		{
			name: "multiple data points for multiple metrics",
			memoryPartition: func() memoryPartition {
				m := memoryPartition{}
				m.metrics.Store("metric1", &metric{
					name: "metric1",
					points: []DataPoint{
						{
							Timestamp: 1,
							Value:     0.1,
						},
						{
							Timestamp: 2,
							Value:     0.2,
						},
					},
					lastIndex: 1,
				})
				m.metrics.Store("metric2", &metric{
					name: "metric2",
					points: []DataPoint{
						{
							Timestamp: 1,
							Value:     0.1,
						},
						{
							Timestamp: 2,
							Value:     0.2,
						},
					},
					lastIndex: 1,
				})
				return m
			}(),
			want: []Row{
				{
					//MetricName: "metric1",
					DataPoint: DataPoint{
						Timestamp: 1,
						Value:     0.1,
					},
				},
				{
					//MetricName: "metric1",
					DataPoint: DataPoint{
						Timestamp: 2,
						Value:     0.2,
					},
				},
				{
					//MetricName: "metric2",
					DataPoint: DataPoint{
						Timestamp: 1,
						Value:     0.1,
					},
				},
				{
					//MetricName: "metric2",
					DataPoint: DataPoint{
						Timestamp: 2,
						Value:     0.2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.memoryPartition.SelectAll()
			assert.Equal(t, tt.want, got)
		})
	}
}