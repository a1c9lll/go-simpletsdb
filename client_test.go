package simpletsdb

import (
	"testing"
	"time"
)

var (
	db *SimpleTSDB
)

func TestMain(t *testing.T) {
	db = NewClient("127.0.0.1", 8981)
}

func TestInsert(t *testing.T) {
	pts := make([]*InsertPointRequest, 10)
	for i := 0; i < len(pts); i++ {
		pts[i] = &InsertPointRequest{
			Metric: "test999",
			Tags: map[string]string{
				"id": "2",
			},
			Point: &Point{
				Value:     25635,
				Timestamp: time.Now().Add(-time.Duration(i) * time.Second).UnixNano(),
			},
		}
	}
	err := db.InsertPoints(pts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryPoints(t *testing.T) {
	ipts := make([]*InsertPointRequest, 5)
	for i := 0; i < len(ipts); i++ {
		ipts[i] = &InsertPointRequest{
			Metric: "test1000",
			Tags: map[string]string{
				"id": "2",
			},
			Point: &Point{
				Value:     24562,
				Timestamp: time.Now().Add(-time.Duration(i) * time.Second).UnixNano(),
			},
		}
	}
	err := db.InsertPoints(ipts)
	if err != nil {
		t.Fatal(err)
	}

	pts, err := db.QueryPoints(&QueryPointsRequest{
		Metric: "test1000",
		Start:  time.Now().Add(-time.Hour).UnixNano(),
		N:      5,
		Tags: map[string]string{
			"id": "2",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(pts) != 5 {
		t.Fatalf("expected 5 points, got %d", len(pts))
	}
}

func TestDeletePoints(t *testing.T) {
	ipts := make([]*InsertPointRequest, 5)
	for i := 0; i < len(ipts); i++ {
		ipts[i] = &InsertPointRequest{
			Metric: "test1001",
			Tags: map[string]string{
				"id": "2",
			},
			Point: &Point{
				Value:     24562,
				Timestamp: time.Now().Add(-time.Duration(i) * time.Second).UnixNano(),
			},
		}
	}

	err := db.InsertPoints(ipts)
	if err != nil {
		t.Fatal(err)
	}

	err = db.DeletePoints(&DeletePointsRequest{
		Metric: "test1001",
		Start:  time.Now().Add(-time.Hour).UnixNano(),
		End:    time.Now().UnixNano(),
		Tags: map[string]string{
			"id": "2",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	pts, err := db.QueryPoints(&QueryPointsRequest{
		Metric: "test1001",
		Start:  time.Now().Add(-time.Hour).UnixNano(),
		N:      5,
		Tags: map[string]string{
			"id": "2",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(pts) != 0 {
		t.Fatalf("expected 0 points, got %d", len(pts))
	}
}

func TestAddDownsampler(t *testing.T) {
	if err := db.AddDownsampler(&Downsampler{
		Metric:    "test10002",
		OutMetric: "test10002_15m_max",
		RunEvery:  "1m",
		Query: &DownsampleQuery{
			Tags: map[string]string{
				"id": "2",
			},
			Window: map[string]interface{}{
				"every": "15m",
			},
			Aggregators: []*AggregatorQuery{
				{Name: "max"},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	dss, err := db.ListDownsamplers()
	if err != nil {
		t.Fatal(err)
	}

	id := int64(-1)
	for _, ds := range dss {
		if ds.OutMetric == "test10002_15m_max" {
			id = ds.ID
		}
	}

	if id == -1 {
		t.Fatal("expected id to not be -1")
	}

	err = db.DeleteDownsampler(&DeleteDownsamplerQuery{
		ID: id,
	})

	if err != nil {
		t.Fatal(err)
	}
}
