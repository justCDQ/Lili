package report

import "testing"

func TestBuild(t *testing.T) {
	input := []Order{
		{ID: "o-2", Status: "paid", AmountCents: 750},
		{ID: "o-1", Status: "paid", AmountCents: 1250},
		{ID: "o-3", Status: "cancelled", AmountCents: 500},
	}
	got, err := Build(input)
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalCents != 2000 || len(got.Paid) != 2 || got.Paid[0].ID != "o-1" {
		t.Fatalf("Build()=%+v", got)
	}
}

func TestBuildRejectsDuplicate(t *testing.T) {
	_, err := Build([]Order{{ID: "o-1", Status: "paid"}, {ID: "o-1", Status: "paid"}})
	if err == nil {
		t.Fatal("Build() error=nil, want duplicate failure")
	}
}
