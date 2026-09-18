package handlers

import (
	"testing"
	"time"
)

func TestLast7DaysWindow(t *testing.T) {
	// 无论输入时刻处于哪个时区，窗口始终按东八区自然日界定，长度恒为 7 天。
	cases := []struct {
		name      string
		now       string
		wantStart string
		wantEnd   string
	}{
		{"CST 零点", "2026-09-17T00:00:00+08:00", "2026-09-11T00:00:00+08:00", "2026-09-18T00:00:00+08:00"},
		{"CST 深夜", "2026-09-17T23:59:59+08:00", "2026-09-11T00:00:00+08:00", "2026-09-18T00:00:00+08:00"},
		{"UTC 零点对应北京早晨", "2026-09-17T00:30:00Z", "2026-09-11T00:00:00+08:00", "2026-09-18T00:00:00+08:00"},
		{"UTC16:00 恰为北京次日零点", "2026-09-16T16:00:00Z", "2026-09-11T00:00:00+08:00", "2026-09-18T00:00:00+08:00"},
		{"负偏移时区输入", "2026-09-17T20:00:00-05:00", "2026-09-12T00:00:00+08:00", "2026-09-19T00:00:00+08:00"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			now, _ := time.Parse(time.RFC3339, tc.now)
			start, end := last7DaysWindow(now)
			if got := start.Format(time.RFC3339); got != tc.wantStart {
				t.Errorf("start = %s, want %s", got, tc.wantStart)
			}
			if got := end.Format(time.RFC3339); got != tc.wantEnd {
				t.Errorf("end = %s, want %s", got, tc.wantEnd)
			}
			if d := end.Sub(start); d != 7*24*time.Hour {
				t.Errorf("window duration = %v, want 168h0m0s", d)
			}
		})
	}
}
