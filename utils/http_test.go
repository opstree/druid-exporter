package utils

import (
	"reflect"
	"testing"
)

func TestGetHealth(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHealth(tt.args.url); got != tt.want {
				t.Errorf("GetHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResponse(t *testing.T) {
	type args struct {
		url       string
		queryType string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResponse(tt.args.url, tt.args.queryType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
