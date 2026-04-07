package dzce

import (
	"testing"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/lintkit/linting"
)

func TestRegisterLintRulesNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRules(nil); err != ErrNilLintRuleRegistrar {
		t.Fatalf(
			"RegisterLintRules(nil) error=%v, want %v",
			err,
			ErrNilLintRuleRegistrar,
		)
	}
}

func TestRegisterLintRules(t *testing.T) {
	t.Parallel()

	engine := linting.NewEngine()
	if err := RegisterLintRules(engine); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	if len(engine.Rules()) == 0 {
		t.Fatal("RegisterLintRules() registered 0 rules, want >0")
	}
}

func TestAttachLintDiagnostics(t *testing.T) {
	t.Parallel()

	engine := linting.NewEngine()
	if err := RegisterLintRules(engine); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	runContext := lint.RunContext{
		TargetPath: "globals.xml",
		TargetKind: lint.FileKind(KindGlobals),
	}
	AttachLintDiagnostics(&runContext, []lint.Diagnostic{
		{
			Code:    lint.ApplyCodePrefix("DZCE", CodeGlobalsInvalidTypeTag),
			Message: "bad type",
		},
	})

	result, err := engine.RunDefault(runContext, nil)
	if err != nil {
		t.Fatalf("engine.RunDefault() error: %v", err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf(
			"engine.RunDefault() diagnostics=%d, want 1",
			len(result.Diagnostics),
		)
	}
}
