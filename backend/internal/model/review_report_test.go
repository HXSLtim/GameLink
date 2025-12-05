package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReviewReportStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status ReviewReportStatus
		want   bool
	}{
		{
			name:   "pending status is valid",
			status: ReviewReportStatusPending,
			want:   true,
		},
		{
			name:   "approved status is valid",
			status: ReviewReportStatusApproved,
			want:   true,
		},
		{
			name:   "rejected status is valid",
			status: ReviewReportStatusRejected,
			want:   true,
		},
		{
			name:   "invalid status",
			status: ReviewReportStatus("invalid"),
			want:   false,
		},
		{
			name:   "empty status",
			status: ReviewReportStatus(""),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Valid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReviewReport_Creation(t *testing.T) {
	report := &ReviewReport{
		ReviewID:   1001,
		ReporterID: 2001,
		Reason:     "评价内容包含不实信息",
		Evidence:   "https://example.com/evidence.jpg",
		Status:     ReviewReportStatusPending,
	}

	assert.Equal(t, uint64(1001), report.ReviewID)
	assert.Equal(t, uint64(2001), report.ReporterID)
	assert.Equal(t, "评价内容包含不实信息", report.Reason)
	assert.Equal(t, "https://example.com/evidence.jpg", report.Evidence)
	assert.Equal(t, ReviewReportStatusPending, report.Status)
	assert.Nil(t, report.HandledBy)
	assert.Nil(t, report.HandledAt)
}
