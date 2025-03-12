package models

import "testing"

func TestDownloadTarget_DownloadSub(t *testing.T) {
	type fields struct {
		Source     TargetType
		Group      string
		Number     string
		Name       string
		localPath  string
		localFiles []string
		sep        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{"non exist files", fields{TargetDmm, "h_1674onez", "165", "onez-154", "test/", []string{}, "-"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DownloadTarget{
				Source:     tt.fields.Source,
				Group:      tt.fields.Group,
				Number:     tt.fields.Number,
				Name:       tt.fields.Name,
				localPath:  tt.fields.localPath,
				localFiles: tt.fields.localFiles,
				sep:        tt.fields.sep,
			}
			if err := d.DownloadSub(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadSub() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
