package dzce

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeLintEconomyCoreTreeMissingInclude(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	corePath := filepath.Join(root, "cfgeconomycore.xml")

	writeTestFile(
		t,
		corePath,
		`<?xml version="1.0"?>
<economycore>
	<ce folder="db">
		<file name="missing.xml" type="types"/>
	</ce>
</economycore>`,
	)

	diagnostics := AnalyzeLintEconomyCoreTree(corePath)
	if !hasCode(diagnostics, CodeMergeMissingIncludeFile) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeMergeMissingIncludeFile))
	}
}

func TestAnalyzeLintEconomyCoreTreeCycle(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	a := filepath.Join(root, "a.xml")
	b := filepath.Join(root, "b.xml")

	writeTestFile(
		t,
		a,
		`<?xml version="1.0"?>
<economycore>
	<ce folder=".">
		<file name="b.xml" type="economycore"/>
	</ce>
</economycore>`,
	)
	writeTestFile(
		t,
		b,
		`<?xml version="1.0"?>
<economycore>
	<ce folder=".">
		<file name="a.xml" type="economycore"/>
	</ce>
</economycore>`,
	)

	diagnostics := AnalyzeLintEconomyCoreTree(a)
	if !hasCode(diagnostics, CodeMergeIncludeCycle) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeMergeIncludeCycle))
	}
}

func TestAnalyzeLintEconomyCoreTreeCrossRefsAndOverrides(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	corePath := filepath.Join(root, "cfgeconomycore.xml")
	dbPath := filepath.Join(root, "db")
	if err := os.MkdirAll(dbPath, 0o750); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}

	writeTestFile(
		t,
		corePath,
		`<?xml version="1.0"?>
<economycore>
	<ce folder="db">
		<file name="types.xml" type="types"/>
		<file name="types_mod.xml" type="types"/>
		<file name="events.xml" type="events"/>
	</ce>
</economycore>`,
	)
	writeTestFile(
		t,
		filepath.Join(dbPath, "types.xml"),
		`<?xml version="1.0"?>
<types>
	<type name="Rag"><nominal>1</nominal></type>
</types>`,
	)
	writeTestFile(
		t,
		filepath.Join(dbPath, "types_mod.xml"),
		`<?xml version="1.0"?>
<types>
	<type name="Rag"><nominal>2</nominal></type>
</types>`,
	)
	writeTestFile(
		t,
		filepath.Join(dbPath, "events.xml"),
		`<?xml version="1.0"?>
<events>
	<event name="A">
		<secondary>B</secondary>
		<children><child type="UnknownType"/></children>
	</event>
</events>`,
	)

	diagnostics := AnalyzeLintEconomyCoreTree(corePath)
	if !hasCode(diagnostics, CodeMergeDuplicateTypeOverride) {
		t.Fatalf(
			"expected %s code in diagnostics",
			codeToken(CodeMergeDuplicateTypeOverride),
		)
	}
	if !hasCode(diagnostics, CodeCrossRefMissingEvent) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeCrossRefMissingEvent))
	}
	if !hasCode(diagnostics, CodeCrossRefMissingType) {
		t.Fatalf("expected %s code in diagnostics", codeToken(CodeCrossRefMissingType))
	}
}

// writeTestFile writes UTF-8 test fixture to target path.
func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) error: %v", path, err)
	}
}
