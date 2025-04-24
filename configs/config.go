package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// DBConfig represents database connection settings
type DBConfig struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DBName   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
	Params   string `yaml:"params"`
}

// Tables represents database table and field configurations
type Tables struct {
	Songs struct {
		Name   string `yaml:"name"`
		Fields struct {
			ID            string `yaml:"id"`
			Name          string `yaml:"name"`
			Artist        string `yaml:"artist"`
			Fingerprinted string `yaml:"fingerprinted"`
			FileSHA1      string `yaml:"file_sha1"`
			TotalHashes   string `yaml:"total_hashes"`
		} `yaml:"fields"`
	} `yaml:"songs"`

	Fingerprints struct {
		Name   string `yaml:"name"`
		Fields struct {
			Hash   string `yaml:"hash"`
			Offset string `yaml:"offset"`
		} `yaml:"fields"`
	} `yaml:"fingerprints"`
}

// Config represents the main application configuration
type Config struct {
	Config struct {
		Name                 string  `yaml:"name"`
		Version              string  `yaml:"version"`
		ConnectivityMask     int     `yaml:"connectivity_mask"`
		SamplingRate         int     `yaml:"sampling_rate"`
		FFTWindowSize        int     `yaml:"fft_window_size"`
		OverlapRatio         float64 `yaml:"overlap_ratio"`
		FanValue             int     `yaml:"fan_value"`
		AmplitudeMin         int     `yaml:"amplitude_min"`
		PeakNeighborhoodSize int     `yaml:"peak_neighborhood_size"`
		MinHashTimeDelta     int     `yaml:"min_hash_time_delta"`
		MaxHashTimeDelta     int     `yaml:"max_hash_time_delta"`
		PeakSort             bool    `yaml:"peak_sort"`
		FingerprintReduction int     `yaml:"fingerprint_reduction"`
		FingerprintLimit     int     `yaml:"fingerprint_limit"`
	} `yaml:"config"`

	Recognition struct {
		TopResults int `yaml:"top_results"`
	} `yaml:"recognition"`

	Database DBConfig `yaml:"database"`
	Tables   Tables   `yaml:"tables"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
