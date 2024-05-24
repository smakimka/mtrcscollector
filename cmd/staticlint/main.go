package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	"github.com/smakimka/mtrcscollector/internal/staticlint"
)

func main() {
	var checks []*analysis.Analyzer

	randomAnalyzersCount := 0
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		} else if randomAnalyzersCount < 5 {
			randomAnalyzersCount += 1
			checks = append(checks, v.Analyzer)
		}
	}

	checks = append(checks, staticlint.MainExitAnalyzer)
	checks = append(checks, appends.Analyzer)
	checks = append(checks, asmdecl.Analyzer)
	checks = append(checks, assign.Analyzer)
	checks = append(checks, atomic.Analyzer)
	checks = append(checks, atomicalign.Analyzer)
	checks = append(checks, bools.Analyzer)
	checks = append(checks, buildssa.Analyzer)
	checks = append(checks, buildtag.Analyzer)
	checks = append(checks, cgocall.Analyzer)
	checks = append(checks, copylock.Analyzer)
	checks = append(checks, composite.Analyzer)
	checks = append(checks, ctrlflow.Analyzer)
	checks = append(checks, deepequalerrors.Analyzer)
	checks = append(checks, defers.Analyzer)
	checks = append(checks, directive.Analyzer)
	checks = append(checks, errorsas.Analyzer)
	checks = append(checks, fieldalignment.Analyzer)
	checks = append(checks, findcall.Analyzer)
	checks = append(checks, framepointer.Analyzer)
	checks = append(checks, httpmux.Analyzer)
	checks = append(checks, httpresponse.Analyzer)
	checks = append(checks, ifaceassert.Analyzer)
	checks = append(checks, loopclosure.Analyzer)
	checks = append(checks, lostcancel.Analyzer)
	checks = append(checks, nilfunc.Analyzer)
	checks = append(checks, nilness.Analyzer)
	checks = append(checks, pkgfact.Analyzer)
	checks = append(checks, printf.Analyzer)
	checks = append(checks, reflectvaluecompare.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, shift.Analyzer)
	checks = append(checks, sigchanyzer.Analyzer)
	checks = append(checks, slog.Analyzer)
	checks = append(checks, sortslice.Analyzer)
	checks = append(checks, stdmethods.Analyzer)
	checks = append(checks, stdversion.Analyzer)
	checks = append(checks, stringintconv.Analyzer)
	checks = append(checks, structtag.Analyzer)
	checks = append(checks, testinggoroutine.Analyzer)
	checks = append(checks, tests.Analyzer)
	checks = append(checks, timeformat.Analyzer)
	checks = append(checks, unmarshal.Analyzer)
	checks = append(checks, unreachable.Analyzer)
	checks = append(checks, unsafeptr.Analyzer)
	checks = append(checks, unusedresult.Analyzer)
	checks = append(checks, unusedwrite.Analyzer)
	checks = append(checks, usesgenerics.Analyzer)

	multichecker.Main(
		checks...,
	)
}
