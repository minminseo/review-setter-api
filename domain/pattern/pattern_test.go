package pattern

import "testing"

var Debug bool = false

func TestValidateSteps(t *testing.T) {
	if Debug {
		t.Skip("スキップさせる")
	}

	s1 := &PatternStep{PatternStepID: "1", UserID: "u1", PatternID: "p1", StepNumber: 1, IntervalDays: 1}
	s2 := &PatternStep{PatternStepID: "2", UserID: "u1", PatternID: "p1", StepNumber: 2, IntervalDays: 2}
	s3 := &PatternStep{PatternStepID: "3", UserID: "u1", PatternID: "p1", StepNumber: 1, IntervalDays: 3} // StepNumber 重複
	s4 := &PatternStep{PatternStepID: "4", UserID: "u1", PatternID: "p1", StepNumber: 3, IntervalDays: 2} // IntervalDays 重複

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
			err := ValidateSteps(tt.args)
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
