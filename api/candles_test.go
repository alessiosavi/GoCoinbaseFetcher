package api

import "testing"

func TestGetCandles(t *testing.T) {
	type args struct {
		pair        string
		startDate   string
		granularity int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok1",
			args: args{
				pair:        "BTC-EUR",
				startDate:   "2021-04-21",
				granularity: 60,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetCandles(tt.args.pair, tt.args.startDate, tt.args.granularity); (err != nil) != tt.wantErr {
				t.Errorf("GetCandles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
