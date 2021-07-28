package simpletsdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Creates a new client
func NewClient(host string, port int) *SimpleTSDB {
	client := &http.Client{}
	return &SimpleTSDB{
		Host:   host,
		Port:   port,
		client: client,
	}
}

func (p *pointReader) Read(b []byte) (int, error) {
	offset := 0
	for ; p.i < len(p.reqs); p.i++ {
		bs := p.reqs[p.i].toLineProtocol()
		if offset+len(bs)+1 >= len(b) {
			return offset, nil
		}
		copy(b[offset:], bs)
		offset += len(bs)
		b[offset] = '\n'
		offset++
	}
	return offset, io.EOF
}

// Inserts points
func (db *SimpleTSDB) InsertPoints(points []*InsertPointRequest) error {
	url := fmt.Sprintf("http://%s:%d/insert_points", db.Host, db.Port)

	reader := &pointReader{reqs: points}
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/x.simpletsdb.points")

	if err != nil {
		return err
	}
	resp, err := db.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Queries points in a time range
func (db *SimpleTSDB) QueryPoints(query *QueryPointsRequest) ([]*Point, error) {
	url := fmt.Sprintf("http://%s:%d/query_points", db.Host, db.Port)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, buf)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}
	resp, err := db.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pts := []*Point{}
	if err := json.NewDecoder(resp.Body).Decode(&pts); err != nil {
		return nil, err
	}

	return pts, nil
}

// Deletes points in a time range
func (db *SimpleTSDB) DeletePoints(request *DeletePointsRequest) error {
	url := fmt.Sprintf("http://%s:%d/delete_points", db.Host, db.Port)

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(request); err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, buf)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}
	resp, err := db.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
