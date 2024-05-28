package net

import (
	"net"
	"reflect"
	"testing"
)

func TestGenerateMask(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid CIDR",
			cidr:    "192.168.10.0/24",
			want:    "255.255.255.0",
			wantErr: false,
		},
		{
			name:    "Invalid CIDR",
			cidr:    "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateMask(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateMask_Validity(t *testing.T) {
	cidr := "192.168.10.0/24"
	want := "255.255.255.0"

	got, err := GenerateMask(cidr)
	if err != nil {
		t.Errorf("GenerateMask() error = %v, wantErr %v", err, false)
		return
	}

	wantIP := net.ParseIP(want)
	gotIP := net.ParseIP(got)

	if !reflect.DeepEqual(wantIP, gotIP) {
		t.Errorf("GenerateMask() = %v, want %v", got, want)
	}
}
