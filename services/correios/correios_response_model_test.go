package correios

import (
	"os"
	"testing"
)

func TestParseResponseModel(t *testing.T) {
	body, err := os.ReadFile("correios_response_model.json")
	if err != nil {
		t.Errorf("Error reading file: %s", err)
		return
	}

	t.Run("default", func(t *testing.T) {
		gotCrm, err := ParseResponseModel(body)
		if err != nil {
			t.Errorf("ParseResponseModel() error = %v", err)
			return
		}
		if gotCrm == nil {
			t.Errorf("ParseResponseModel() gotCrm = nil")
		}
		if len(gotCrm.Objetos) != 1 {
			t.Errorf("ParseResponseModel() len(gotCrm.Objetos) = %v", len(gotCrm.Objetos))
		}
		if len(gotCrm.Objetos[0].Eventos) != 8 {
			t.Errorf("ParseResponseModel() len(gotCrm.Objetos[0].Eventos) = %v", len(gotCrm.Objetos[0].Eventos))
		}
	})

}
