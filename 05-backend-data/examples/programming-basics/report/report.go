package report

import (
	"fmt"
	"sort"
)

type Order struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	AmountCents int64  `json:"amount_cents"`
}

type Report struct {
	Paid       []Order `json:"paid"`
	TotalCents int64   `json:"total_cents"`
}

func Build(orders []Order) (Report, error) {
	report := Report{Paid: make([]Order, 0, len(orders))}
	seen := make(map[string]struct{}, len(orders))
	for i, order := range orders {
		if order.ID == "" {
			return Report{}, fmt.Errorf("order[%d]: empty id", i)
		}
		if _, exists := seen[order.ID]; exists {
			return Report{}, fmt.Errorf("order[%d]: duplicate id %q", i, order.ID)
		}
		seen[order.ID] = struct{}{}
		if order.AmountCents < 0 {
			return Report{}, fmt.Errorf("order[%d]: negative amount", i)
		}
		switch order.Status {
		case "paid":
			report.Paid = append(report.Paid, order)
			report.TotalCents += order.AmountCents
		case "cancelled":
		default:
			return Report{}, fmt.Errorf("order[%d]: invalid status %q", i, order.Status)
		}
	}
	sort.Slice(report.Paid, func(i, j int) bool { return report.Paid[i].ID < report.Paid[j].ID })
	return report, nil
}
