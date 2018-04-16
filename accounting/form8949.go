package accounting

import (
	"encoding/csv"
	"os"
	"sort"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type Form8949 struct {
	LongHolds  []Form8949LineItem
	ShortHolds []Form8949LineItem
}

func (form *Form8949) sort() {
	sort.Slice(form.ShortHolds, func(i, j int) bool {
		return form.ShortHolds[i].DateAcquired.Before(form.ShortHolds[j].DateAcquired) ||
			form.ShortHolds[i].DateAcquired.Equal(form.ShortHolds[j].DateAcquired)
	})
	sort.Slice(form.LongHolds, func(i, j int) bool {
		return form.LongHolds[i].DateAcquired.Before(form.LongHolds[j].DateAcquired) ||
			form.LongHolds[i].DateAcquired.Equal(form.LongHolds[j].DateAcquired)
	})
	/*
		sort.Slice(form.ShortHolds, func(i, j int) bool {
			return len(form.ShortHolds[i].Currency) < len(form.ShortHolds[j].Currency)
		})
		sort.Slice(form.LongHolds, func(i, j int) bool {
			return len(form.LongHolds[i].Currency) < len(form.LongHolds[j].Currency)
		})
	*/
}

func (form *Form8949) WriteCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	form.sort()

	emptyRow := []string{}
	header := []string{
		"Description",
		"Date Acquired",
		"Date Sold",
		"Proceeds",
		"Cost Basis",
		"Adjustment Code",
		"Adjustment amount",
		"Gain or loss"}

	writer.Write(emptyRow)
	writer.Write([]string{"Short-Term Holds"})
	writer.Write(header)

	for _, shortHold := range form.ShortHolds {
		err := writer.Write(shortHold.Record())
		if err != nil {
			return err
		}
	}

	writer.Write(emptyRow)
	writer.Write(emptyRow)
	writer.Write([]string{"Long-Term Holds"})
	writer.Write(header)

	for _, longHold := range form.LongHolds {
		err := writer.Write(longHold.Record())
		if err != nil {
			return err
		}
	}

	return nil
}

type Form8949LineItem struct {
	Currency         string
	Description      string
	DateAcquired     time.Time
	DateSold         time.Time
	Proceeds         string
	CostBasis        string
	AdjustmentCode   string
	AdjustmentAmount string
	GainOrLoss       string
}

func (item *Form8949LineItem) Record() []string {
	return []string{
		item.Description,
		item.DateAcquired.Format(common.TIME_DISPLAY_FORMAT),
		item.DateSold.Format(common.TIME_DISPLAY_FORMAT),
		item.Proceeds,
		item.CostBasis,
		item.AdjustmentCode,
		item.AdjustmentAmount,
		item.GainOrLoss}
}
