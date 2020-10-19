package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"moh.gov.bz/mch/emtct/internal/models"
)

func TestApp_RetrievePatient(t *testing.T) {
	patientId := "1111120"
	r := RegisterHandlers()

	req, err := http.NewRequest("GET", fmt.Sprintf("/patient/%s", patientId), nil)
	if err != nil {
		t.Fatalf("error creating request: %+v", err)
	}
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImVhLUFIZi0wN0FaeGVCNzhSVkVoUyJ9.eyJuaWNrbmFtZSI6InJvYmVydG8uZ3VlcnJhIiwibmFtZSI6IlJvYmVydG8gR3VlcnJhIiwicGljdHVyZSI6Imh0dHBzOi8vcy5ncmF2YXRhci5jb20vYXZhdGFyLzQyYWY1ZWFjMjdmMzBiNTE4NzRiODQwMGY4YWIxOGZiP3M9NDgwJnI9cGcmZD1odHRwcyUzQSUyRiUyRmNkbi5hdXRoMC5jb20lMkZhdmF0YXJzJTJGcm8ucG5nIiwidXBkYXRlZF9hdCI6IjIwMjAtMTAtMThUMTU6NDk6MzguNjI2WiIsImVtYWlsIjoicm9iZXJ0by5ndWVycmFAb3BlbnN0ZXAubmV0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vZW10Y3QtZGV2LnVzLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1ZjgwZmRhMzg4YTI1YTAwNmJiZDlmMzYiLCJhdWQiOiJrNDZoZmJCVURzT2FQZ05VOUlsVWQ3aG9XSjVLdTBFQiIsImlhdCI6MTYwMzAzNjM0MywiZXhwIjoxNjAzMDcyMzQzLCJub25jZSI6ImNWZExlVUl5TkhCblpsTlJORGxsUkdwSlkzZGpUSEpXVkZSM01qSkJZVTVuWVVGVFJYSmlkbkkxZFE9PSJ9.u6nXzKWe7Z54a0dAog72xLjcv0z4DvrS8psISM6izbBH4W46AcGtNnwEUDF_UW29AcVes6rjt8v8Gkf_dGKnhA8ACDTfNEDP1DMXWGi7IM1V8HZzF9EfVvqjGff1dp2QRsVDGnKpTcDW2tMud5PnjyJOLF_fK8K3K0dbH3LvtH1qYX5gAfa0I0ROY9rbnv7WxFKhKHKhZSrNerFjj0MthSdw5rBVGIrC0mr4GPpiCrssep2ye4KtCnsfadAjaRZFm8gar0lpEjpu4-G7IfIuqQr8oTDHx3F-yfP0a6XgMP-B9z0ItV69HadBXHoAIIJzpDel9Resme2TuCSAG_-stQ")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code error, want 200, got %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var patient models.Patient
	_ = json.Unmarshal(body, &patient)
	t.Logf("resp: %+v", patient)
	if patient.Id != patientId {
		t.Errorf("want: %s, got: %s", patientId, patient.Id)
	}

}
