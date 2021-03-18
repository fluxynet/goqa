package main

import (
	"reflect"
	"testing"
)

func TestLoadConf(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConf()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
