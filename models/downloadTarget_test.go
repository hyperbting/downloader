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

func TestDownloadTarget_sanitizeName(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{"t0", fields{"濡れてテカってピッタリ密着神スク水跡美しゅり美少女から人妻まで可愛い女子のスクール水着姿をじっとりと堪能！濡れてテカってピッタリ密着神スク水跡美しゅり美少女から人妻まで可愛い女子のスクール水着姿をじっとりと堪能！"}},
		{"t0-1", fields{"濡れてテカってピッタリ密着神スク水跡美しゅり美少女から人妻まで可愛い女子のスクール水着姿をじっとりと堪能濡れてテカってピッタリ密着神スク水跡美しゅり美少女から人妻まで可愛い女子のスクール水着姿をじっとりと堪能 跡美しゅり"}},

		{"t0-2", fields{"t1p4eh189hfguoirehbg0p943hg02387hg07542hg0548g7h0"}},

		{"t1", fields{"t1p4eh189hfguoirehbg0p943hg02387hg07542hg0548g7h0 gfd;kgjfd;lkgjfdk;lgjfd;kgljfd;gkjfdk;gljfd;klgjfd;lgjfd;klgjfd;lkg dfghfdjkgoewfgbkdsjfgsdkjhfgdskhfgkh t123"}},
		{"t1-2", fields{"1234567890 2234567890 3234567890 4234567890 5234567890 6234567890 7234567890 8234567890 9234567890 0234567890 1134567890 t123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DownloadTarget{
				Name: tt.fields.Name,
			}
			d.sanitizeName()

			t.Log(d)
		})
	}
}
