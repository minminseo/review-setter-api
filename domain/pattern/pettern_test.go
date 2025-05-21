package pattern

import "testing"

var Debug bool = false

func TestValidateSteps(t *testing.T) {
	if Debug {
		t.Skip("スキップさせる")
	}

	s1 := &PatternStep{id: "1", patternID: "p1", stepNumber: 1, intervalDays: 1}
	s2 := &PatternStep{id: "2", patternID: "p1", stepNumber: 2, intervalDays: 2}
	s3 := &PatternStep{id: "3", patternID: "p1", stepNumber: 1, intervalDays: 3} // duplicate step number
	s4 := &PatternStep{id: "4", patternID: "p1", stepNumber: 3, intervalDays: 2} // duplicate intervalDays

	tests := []struct {
		name    string
		args    []*PatternStep
		wantErr bool
		errMsg  string
	}{
		// 異常系
		{
			name:    "復習日間隔数が0",
			args:    []*PatternStep{},
			wantErr: true,
			errMsg:  "復習日間隔数は1つ以上指定してください",
		},

		// 正常系
		{
			name:    "復習日間隔数が1つ",
			args:    []*PatternStep{s1},
			wantErr: false,
		},

		// 正常系
		{
			name:    "復習日間隔数が2つ以上かつ昇順",
			args:    []*PatternStep{s1, s2},
			wantErr: false,
		},

		// 異常系
		{
			name:    "順序番号が昇順でない",
			args:    []*PatternStep{s4, s3},
			wantErr: true,
			errMsg:  "順序番号は昇順で指定してください",
		},

		// 異常系
		{
			name:    "復習日間隔数が昇順でない",
			args:    []*PatternStep{s3, s2},
			wantErr: true,
			errMsg:  "復習日間隔数は昇順で指定してください",
		},

		// 異常系
		{
			name:    "順序番号が前の値から+1でない（公差1の等差数列でない）",
			args:    []*PatternStep{s1, s4},
			wantErr: true,
			errMsg:  "順序番号は1ずつ増加して指定してください",
		},

		// 異常系
		{
			name:    "順序番号が重複している",
			args:    []*PatternStep{s1, s3},
			wantErr: true,
			errMsg:  "順序番号は重複してはいけません",
		},

		// 異常系
		{
			name:    "復習日間隔数が重複している",
			args:    []*PatternStep{s1, s4, s2},
			wantErr: true,
			errMsg:  "順序番号は1ずつ増加して指定してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSteps(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("err==nilのため、%qのテストが失敗しました", tt.name)
				}
				if err.Error() != tt.errMsg {
					t.Errorf("予期しないエラー:実際の結果 %q, 期待 %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("予期しないエラー: %v", err)
				}
			}
		})
	}
}
