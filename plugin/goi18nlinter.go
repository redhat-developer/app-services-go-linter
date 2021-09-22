package main

import (
	"github.com/alexal/go-i18n-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

// AnalyzerPlugin analyzer plugin
var AnalyzerPlugin analyzerPlugin

// GetAnalyzers returns all analyzers for a plugin
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}
