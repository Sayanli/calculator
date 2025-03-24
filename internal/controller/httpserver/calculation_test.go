package httpserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sayanli/calculator/internal/entity"
	"github.com/sayanli/calculator/internal/service"
)

type TestCase struct {
	name       string
	input      []entity.Instruction
	expected   *[]entity.Result
	StatusCode int
}

func TestCalculate(t *testing.T) {
	buf := &bytes.Buffer{}
	log := slog.New(slog.NewTextHandler(buf, nil))
	svc := service.NewCalculationService(log, 100)
	router := NewCalculationRouter(svc)
	url := "http://localhost:8080/v1/calculate"

	cases := []TestCase{
		TestCase{
			name: "Addition",
			input: []entity.Instruction{
				{
					Type:  "calc",
					Op:    "+",
					Var:   "x",
					Left:  1,
					Right: 2,
				},
				{
					Type: "print",
					Var:  "x",
				},
			},
			expected: &[]entity.Result{
				{
					Var:   "x",
					Value: 3,
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			name: "StatusBadRequest",
			input: []entity.Instruction{
				{
					Type:  "calc",
					Op:    "fdsf",
					Var:   "x",
					Left:  1,
					Right: 2,
				},
			},
			expected:   nil,
			StatusCode: http.StatusBadRequest,
		},
	}

	for index, item := range cases {
		jsonData, err := json.Marshal(item.input)
		if err != nil {
			t.Fatalf("Error marshaling input: %v", err)
		}
		body := bytes.NewBuffer(jsonData)
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.Calculate(w, req)
		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d", index, w.Code, item.StatusCode)
		}
		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("[%d] Error reading response body: %v", index, err)
		}
		if len(data) == 0 {
			t.Errorf("[%d] Empty response body", index)
		}
	}
}
