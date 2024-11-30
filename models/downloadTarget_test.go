package models

import "testing"

func TestDownloadTarget_TryDownloadDmmMain(t *testing.T) {
	type fields struct {
		Group  string
		Number string
		Name   string
		Source string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{"t0", fields{"dass", "001", "optional", string(TargetDmm)}, false},
		{"t1", fields{"pow", "035", "optional", string(TargetDmm)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DownloadTarget{
				Group:  tt.fields.Group,
				Number: tt.fields.Number,
				Name:   tt.fields.Name,
				Source: tt.fields.Source,
			}
			if err := d.TryDownloadDmmMain(); (err != nil) != tt.wantErr {
				t.Errorf("TryDownloadDmmMain() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("TryDownloadDmmMain() %v", d)
		})
	}
}

func TestDownloadTarget_BuildDmmSubPath(t *testing.T) {
	type fields struct {
		Group  string
		Number string
		Name   string
		Source string
	}
	type args struct {
		sep string
		cnt int
		hd  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"t0", fields{"dass", "001", "optional", string(TargetDmm)}, args{"00", 1, "jp"}, false},
		{"t1", fields{"41zb", "013", "optional", string(TargetDmm)}, args{"00", 1, "jp"}, false},
		{"t2", fields{"51vs", "595", "optional", string(TargetDmm)}, args{"00", 1, "jp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DownloadTarget{
				Group:  tt.fields.Group,
				Number: tt.fields.Number,
				Name:   tt.fields.Name,
				Source: tt.fields.Source,
			}
			if err := d.DownloadSub(""); (err != nil) != tt.wantErr {
				t.Errorf("BuildDmmSubPath() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("BuildDmmSubPath() %v", d)
		})
	}
}
