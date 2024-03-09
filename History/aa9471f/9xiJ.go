package breeze

import (
	"reflect"
	"testing"
)

func TestApificationBreeze_makeBreezeRESTCall(t *testing.T) {
	type fields struct {
		Breeze             *Breeze
		hostname           string
		base64SessionToken string
	}
	type args struct {
		method      string
		endpoint    string
		requestBody map[string]string
		useProxy    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ApificationBreeze{
				Breeze:             tt.fields.Breeze,
				hostname:           tt.fields.hostname,
				base64SessionToken: tt.fields.base64SessionToken,
			}
			got, err := a.executeBreezeAPIRESTRequest(tt.args.method, tt.args.endpoint, tt.args.requestBody, tt.args.useProxy)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApificationBreeze.makeBreezeRESTCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApificationBreeze.makeBreezeRESTCall() = %v, want %v", got, tt.want)
			}
		})
	}
}
