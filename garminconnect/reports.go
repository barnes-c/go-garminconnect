package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// ReportExport identifies a requested report. Poll ReportStatus for ReportID
// after about EstimatedTimeForStatusCheck milliseconds.
type ReportExport struct {
	ReportID                    string `json:"reportId"`
	EstimatedTimeForStatusCheck int    `json:"estimatedTimeForStatusCheck"`
}

// ReportExportStatus is the progress of a requested report. Status is "READY"
// once it can be downloaded; "DATA_NOT_FOUND" means the range had no data.
type ReportExportStatus struct {
	ReportID string `json:"reportId"`
	Progress struct {
		Status          string  `json:"status"`
		ResponseMessage *string `json:"responseMessage"`
	} `json:"progress"`
	CreationTimestampGMT string `json:"creationTimestampGMT"`
}

// RequestReport starts a PDF export of a report category between start and
// end. For "bloodPressure", options must set "reportSubCategory" to "Daily",
// "SevenDays", "FourWeeks" or "Yearly".
func (c *Client) RequestReport(ctx context.Context, category string, start, end time.Time, options map[string]any) (*ReportExport, error) {
	if options == nil {
		options = map[string]any{}
	}
	path := fmt.Sprintf("/reportexport-service/export/init/%s/%s/%s", url.PathEscape(category), date(start), date(end))
	var out ReportExport
	if err := c.post(ctx, path, options, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReportStatus returns the progress of a report requested with RequestReport.
func (c *Client) ReportStatus(ctx context.Context, reportID string) (*ReportExportStatus, error) {
	var out ReportExportStatus
	if err := c.get(ctx, "/reportexport-service/export/reportStatus/"+url.PathEscape(reportID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadReport returns a finished report as a PDF.
func (c *Client) DownloadReport(ctx context.Context, reportID string) ([]byte, error) {
	return c.getBytes(ctx, "/reportexport-service/export/reportContent/"+url.PathEscape(reportID), nil)
}
