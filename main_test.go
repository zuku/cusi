package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	type args struct {
		command int8
		data    string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "listDir",
			args: args{command: 0x03, data: "/flash"},
			want: []byte{0x03, 0x2f, 0x66, 0x6c, 0x61, 0x73, 0x68},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createCommand(tt.args.command, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCommand = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyReceivedContainer(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{data: []byte{0xaa, 0xab, 0xaa, 0x00, 0x00, 0x00, 0x00, 0xab, 0xcc, 0xab}},
			want: true,
		},
		{
			name: "mismatch_head",
			args: args{data: []byte{0xaa, 0xab, 0xac, 0x00, 0x00, 0x00, 0x00, 0xab, 0xcc, 0xab}},
			want: false,
		},
		{
			name: "mismatch_foot",
			args: args{data: []byte{0xaa, 0xab, 0xaa, 0x00, 0x00, 0x00, 0x00, 0xab, 0xac, 0xab}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyReceivedContainer(tt.args.data); got != tt.want {
				t.Errorf("verifyReceivedContainer = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractReceivedData(t *testing.T) {
	type args struct {
		data []byte
	}
	type results struct {
		data    []byte
		iserror bool
	}
	tests := []struct {
		name string
		args args
		want results
	}{
		{
			name: "ok",
			args: args{data: []byte{0xaa, 0xab, 0xaa, 0x00, 0x00, 0x61, 0x70, 0x70, 0x73, 0x2c, 0x00, 0x00, 0xab, 0xcc, 0xab}},
			want: results{data: []byte{0x61, 0x70, 0x70, 0x73, 0x2c}, iserror: false},
		},
		{
			name: "error_response",
			args: args{data: []byte{0xaa, 0xab, 0xaa, 0x05, 0xff, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x00, 0x00, 0xab, 0xcc, 0xab}},
			want: results{data: nil, iserror: true},
		},
		{
			name: "short_response",
			args: args{data: []byte{0xaa, 0xab, 0xaa, 0x00, 0x00}},
			want: results{data: nil, iserror: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractReceivedData(tt.args.data)
			if tt.want.iserror {
				if err == nil {
					t.Errorf("extractReceivedData error is nil, error expected")
				}
			} else {
				if err != nil {
					t.Errorf("extractReceivedData error is not nil, nil expected")
				}
			}
			if !reflect.DeepEqual(got, tt.want.data) {
				t.Errorf("extractReceivedData []byte = %v, want %v", got, tt.want.data)
			}
		})
	}
}

func TestCrc16(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "1234567890",
			args: args{data: []byte("1234567890")},
			want: 0xc20a,
		},
		{
			name: "ABCDEFG",
			args: args{data: []byte("ABCDEFG")},
			want: 0x9e77,
		},
		{
			name: "abcdefg",
			args: args{data: []byte("abcdefg")},
			want: 0xe9c2,
		},
		{
			name: "NULL",
			args: args{data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			want: 0x0770,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crc16(tt.args.data); got != tt.want {
				t.Errorf("crc16 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendCrc(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1234567890",
			args: args{data: []byte("1234567890")},
			want: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0xc2, 0x0a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendCrc(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendCrc = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	type args struct {
		path string
	}
	type results struct {
		path      string
		errortext string
	}
	tests := []struct {
		name string
		args args
		want results
	}{
		{
			name: "normal",
			args: args{path: "path"},
			want: results{path: "/flash/path", errortext: ""},
		},
		{
			name: "root",
			args: args{path: "/"},
			want: results{path: "", errortext: "absolute path is not permitted"},
		},
		{
			name: "root_windows",
			args: args{path: "\\"},
			want: results{path: "", errortext: "absolute path is not permitted"},
		},
		{
			name: "absolute",
			args: args{path: "/root"},
			want: results{path: "", errortext: "absolute path is not permitted"},
		},
		{
			name: "parent",
			args: args{path: "../"},
			want: results{path: "", errortext: "forbidden path"},
		},
		{
			name: "complex",
			args: args{path: "./path/./to/../../other/"},
			want: results{path: "/flash/other", errortext: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizePath(tt.args.path)
			if len(tt.want.errortext) > 0 {
				if !strings.Contains(err.Error(), tt.want.errortext) {
					t.Errorf("normalizePath err = %v, want %v", err, tt.want.errortext)
				}
			} else {
				if err != nil {
					t.Errorf("normalizePath err = %v, want nil", err)
				}
				if got != tt.want.path {
					t.Errorf("normalizePath path = %v, want %v", got, tt.want.path)
				}
			}
		})
	}

}
