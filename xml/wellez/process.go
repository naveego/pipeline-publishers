package wellez

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

func processFile(ctx publisher.Context, transport publisher.DataTransport, tmpDir string, file fileInfo) error {

	xmlFile := strings.Replace(file.FileName, ".zip", ".xml", -1)
	ctx.Logger.Infof("Processing xml file '%s", xmlFile)

	f, err := os.Open(filepath.Join(tmpDir, xmlFile))
	if err != nil {
		return err
	}
	defer f.Close()

	data := Data{}
	err = xml.NewDecoder(f).Decode(&data)
	if err != nil {
		return err
	}

	if err := sendDataPoints(ctx, transport, data.CompletionCost, "CompletionCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.CompletionCostItem, "CompletionCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.DailyOps, "DailyOps", []string{"well_id", "report_date", "job_number"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.DrillingCost, "DrillingCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.DrillingCostItem, "DrillingCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.FacilitiesCost, "FacilitiesCost", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.FacilitiesCostItem, "FacilitiesCostItem", []string{"RowID"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.JobDetails, "JobDetails", []string{"well_id", "job_number"}); err != nil {
		return err
	}
	if err := sendDataPoints(ctx, transport, data.WellInfo, "WellInfo", []string{"well_id"}); err != nil {
		return err
	}

	return nil
}

func toJSONMap(s interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	b, _ := json.Marshal(s)
	d := json.NewDecoder(strings.NewReader(string(b)))
	d.UseNumber()
	d.Decode(&m)
	return m
}

func sendDataPoints(ctx publisher.Context, transport publisher.DataTransport, data interface{}, entity string, keyNames []string) error {

	dataPoints := []pipeline.DataPoint{}
	count := 0

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			datum := s.Index(i).Interface()
			dp := pipeline.DataPoint{
				Repository: ctx.PublisherInstance.Repository,
				Source:     ctx.PublisherInstance.SafeName,
				Entity:     entity,
				Action:     "upsert",
				KeyNames:   keyNames,
				Data:       toJSONMap(datum),
			}

			dataPoints = append(dataPoints, dp)
			count++

			if count == 100 {
				if err := transport.Send(dataPoints); err != nil {
					return err
				}

				count = 0
				dataPoints = []pipeline.DataPoint{}
			}
		}
	}

	if len(dataPoints) > 0 {
		if err := transport.Send(dataPoints); err != nil {
			return err
		}
	}

	return nil
}
