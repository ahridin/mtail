// Copyright 2019 Google Inc. All Rights Reserved.
// This file is available under the Apache license.
// +build integration

package mtail_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/mtail/internal/mtail"
	"github.com/google/mtail/internal/testutil"
)

func TestNewProg(t *testing.T) {
	tmpDir, rmTmpDir := testutil.TestTempDir(t)
	defer rmTmpDir()

	logDir := path.Join(tmpDir, "logs")
	progDir := path.Join(tmpDir, "progs")
	err := os.Mkdir(logDir, 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(progDir, 0700)
	if err != nil {
		t.Fatal(err)
	}

	m, stopM := mtail.TestStartServer(t, 0, false, mtail.ProgramPath(progDir), mtail.LogPathPatterns(logDir+"/*"))
	defer stopM()

	startProgLoadsTotal := mtail.TestGetMetric(t, m.Addr(), "prog_loads_total").(map[string]interface{})

	testutil.TestOpenFile(t, progDir+"/nocode.mtail")
	time.Sleep(time.Second)

	progLoadsTotal := mtail.TestGetMetric(t, m.Addr(), "prog_loads_total").(map[string]interface{})

	mtail.ExpectMetricDelta(t, progLoadsTotal["nocode.mtail"], startProgLoadsTotal["nocode.mtail"], 1)

}
