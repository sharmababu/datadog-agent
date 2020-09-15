// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// +build functionaltests

package tests

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/security/rules"
)

func TestProcess(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	ruleDef := &rules.RuleDefinition{
		ID:         "test_rule",
		Expression: fmt.Sprintf(`process.user == "%s" && process.name == "%s" && open.filename == "/etc/hosts"`, currentUser.Name, path.Base(executable)),
	}

	test, err := newTestModule(nil, []*rules.RuleDefinition{ruleDef}, testOpts{enableFilters: true})
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	f, err := os.Open("/etc/hosts")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, rule, err := test.GetEvent()
	if err != nil {
		t.Error(err)
	} else {
		if rule.ID != "test_rule" {
			t.Errorf("expected rule 'test-rule' to be triggered, got %s", rule.ID)
		}
	}
}

func TestProcessFilename(t *testing.T) {
	ruleDef := &rules.RuleDefinition{
		ID:         "test_rule",
		Expression: fmt.Sprintf(`open.filename == "/etc/hosts"`),
	}

	test, err := newTestModule(nil, []*rules.RuleDefinition{ruleDef}, testOpts{enableFilters: true})
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	f, err := os.Open("/etc/hosts")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	event, _, err := test.GetEvent()
	if err != nil {
		t.Error(err)
	} else {
		if filename, _ := event.GetFieldValue("process.filename"); filename == "" {
			t.Error("should get a valid filename")
		}
	}
}
