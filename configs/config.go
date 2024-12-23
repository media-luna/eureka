package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DBConfig struct {
    Type     string `yaml:"type"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Host     string `yaml:"host"`
    DBName   string `yaml:"db_name"`
    Port     int    `yaml:"port"`
    Params   string `yaml:"params"`
    Supported []string `yaml:"supported"`
}

type Tables struct{
    Songs struct {
        Name   string `yaml:"name"`
        Fields struct {
            ID          string `yaml:"id"`
            Name        string `yaml:"name"`
            Fingerprinted string `yaml:"fingerprinted"`
            FileSHA1    string `yaml:"file_sha1"`
            TotalHashes string `yaml:"total_hashes"`
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

type Config struct {
    EurekaResponse struct {
        SongID              string `yaml:"song_id"`
        SongName            string `yaml:"song_name"`
        Results             string `yaml:"results"`
        HashesMatched       string `yaml:"hashes_matched"`
        FingerprintedHashes string `yaml:"fingerprinted_hashes"`
        FingerprintedConfidence string `yaml:"fingerprinted_confidence"`
        InputHashes         string `yaml:"input_hashes"`
        InputConfidence     string `yaml:"input_confidence"`
        TotalTime           string `yaml:"total_time"`
        FingerprintTime     string `yaml:"fingerprint_time"`
        QueryTime           string `yaml:"query_time"`
        AlignTime           string `yaml:"align_time"`
        Offset              string `yaml:"offset"`
        OffsetSeconds       string `yaml:"offset_seconds"`
    } `yaml:"eureka_response"`
    
    Database DBConfig `yaml:"database"`
    Tables Tables `yaml:"tables"`

    // SQLTemplates map[string]map[string]string `yaml:"sql_templates"`
    SQLTemplates struct {
        MySQL    string `yaml:"mysql"`
        Postgres string `yaml:"postgres"`
        Template struct {
            CreateSongsTable string `yaml:"create_songs_table"`
            CreateFingerprintsTable string `yaml:"create_fingerprints_table"`
            DeleteUnfingerprinted string `yaml:"delete_unfingerprinted"`
        } `yaml:"template"`
    } `yaml:"sql_templates"`

    Config struct {
        Name                string  `yaml:"name"`
        Version             string  `yaml:"version"`
        ConnectivityMask    int     `yaml:"connectivity_mask"`
        SamplingRate        int     `yaml:"sampling_rate"`
        FFTWindowSize       int     `yaml:"fft_window_size"`
        OverlapRatio        float64 `yaml:"overlap_ratio"`
        FanValue            int     `yaml:"fan_value"`
        AmplitudeMin        int     `yaml:"amplitude_min"`
        PeakNeighborhoodSize int    `yaml:"peak_neighborhood_size"`
        MinHashTimeDelta    int     `yaml:"min_hash_time_delta"`
        MaxHashTimeDelta    int     `yaml:"max_hash_time_delta"`
        PeakSort            bool    `yaml:"peak_sort"`
        FingerprintReduction int    `yaml:"fingerprint_reduction"`
    } `yaml:"config"`

    Recognition struct {
        TopResults int `yaml:"top_results"`
    } `yaml:"recognition"`
}



func LoadConfig(filePath string) (*Config, error) {
	// Open the YAML file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the YAML file
	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
