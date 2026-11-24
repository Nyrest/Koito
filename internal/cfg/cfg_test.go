package cfg

import "testing"

func TestLoadConfigBackfillMediaOnStartup(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "unset", want: true},
		{name: "enabled", value: "true", want: true},
		{name: "disabled", value: "false", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := map[string]string{BACKFILL_MEDIA_ON_STARTUP_ENV: tt.value}
			cfg, err := loadConfig(func(key string) string { return env[key] }, "test")
			if err != nil {
				t.Fatal(err)
			}
			if cfg.backfillMediaOnStartup != tt.want {
				t.Errorf("backfillMediaOnStartup = %t, want %t", cfg.backfillMediaOnStartup, tt.want)
			}
		})
	}
}
