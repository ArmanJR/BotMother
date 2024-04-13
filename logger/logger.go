package logger

import (
	"BotMother/config"
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var log *zap.Logger

func Init() {
	// Configure the Elasticsearch client with the local URL
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Configs.ElasticSearchHost, // Local Elasticsearch URL
		},
	}

	// Initialize the Elasticsearch client with the configuration
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Error creating the Elasticsearch client: %s", err))
	}

	// Define a custom zapcore.Core to integrate with Elasticsearch
	core := newElasticsearchCore(zapcore.InfoLevel, zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), es)

	// Create the global logger with the Elasticsearch core and the standard error output
	log = zap.New(zapcore.NewTee(core, zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.Lock(os.Stderr), zapcore.DebugLevel)))
}

// Logger returns the global logger
func Logger() *zap.Logger {
	return log
}

// newElasticsearchCore creates a new zapcore.Core that writes logs to Elasticsearch
func newElasticsearchCore(enab zapcore.LevelEnabler, enc zapcore.Encoder, es *elasticsearch.Client) zapcore.Core {
	return &esCore{
		LevelEnabler: enab,
		enc:          enc,
		esClient:     es,
	}
}

// esCore is a custom zapcore.Core implementation for Elasticsearch
type esCore struct {
	zapcore.LevelEnabler
	enc      zapcore.Encoder
	esClient *elasticsearch.Client
}

func (ec *esCore) With(fields []zapcore.Field) zapcore.Core {
	clone := ec.clone()
	for _, field := range fields {
		field.AddTo(clone.enc)
	}
	return clone
}

func (ec *esCore) clone() *esCore {
	return &esCore{
		LevelEnabler: ec.LevelEnabler,
		enc:          ec.enc.Clone(),
		esClient:     ec.esClient,
	}
}

func (ec *esCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ec.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, ec)
	}
	return checkedEntry
}

func (ec *esCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buffer, err := ec.enc.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	// Send the log entry to Elasticsearch
	indexName := fmt.Sprintf("logs-%s", time.Now().Format("2006.01.02")) // daily index format
	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    bytes.NewReader(buffer.Bytes()),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), ec.esClient)
	if err != nil {
		return fmt.Errorf("error sending log to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}
	return nil
}

func (ec *esCore) Sync() error {
	return nil
}
