package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/providers"
)

func Run(ctx context.Context, args []string) (int, error) {
	if len(args) >= 1 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		printHelp()
		return ExitOK, nil
	}

	if len(args) >= 1 && (args[0] == "--version" || args[0] == "version") {
		fmt.Printf("weave %s\n", Version())
		return ExitOK, nil
	}

	if len(args) >= 1 && args[0] == "doctor" {
		return runDoctor(ctx, args[1:])
	}

	if len(args) >= 1 && args[0] == "migrate" {
		return runMigrate(ctx, args[1:])
	}

	if len(args) == 0 {
		firstRun := !configExists(resolveWorkdir())
		h := NewDefaultForgeHandler()
		code, result, err := h.RunWithOptions(ctx, false)
		if err != nil {
			return code, err
		}
		printForgeSummary(result.ServiceResult)
		if firstRun {
			printFirstRunGuidance()
		}
		return code, nil
	}

	if args[0] == "forge" {
		dryRun, err := parseDryRunOnly(args[1:], "forge")
		if err != nil {
			return exitCodeForError(err), err
		}
		firstRun := !dryRun && !configExists(resolveWorkdir())
		h := NewDefaultForgeHandler()
		code, result, err := h.RunWithOptions(ctx, dryRun)
		if err != nil {
			return code, err
		}
		printForgeSummary(result.ServiceResult)
		if firstRun {
			printFirstRunGuidance()
		}
		return code, nil
	}

	if len(args) >= 3 && args[0] == "skill" && args[1] == "add" {
		return runAssetAdd(ctx, assetKindSkill, args[2], args[3:])
	}

	if len(args) >= 3 && args[0] == "command" && args[1] == "add" {
		return runAssetAdd(ctx, assetKindCommand, args[2], args[3:])
	}

	if len(args) >= 2 && args[0] == "provider" {
		return runProvider(ctx, args[1:])
	}

	return ExitRuntimeError, fmt.Errorf("unsupported command")
}

func runMigrate(ctx context.Context, args []string) (int, error) {
	dryRun, err := parseDryRunOnly(args, "migrate")
	if err != nil {
		return exitCodeForError(err), err
	}

	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return exitCodeForError(app.WrapInvalidConfig(err)), app.WrapInvalidConfig(err)
	}

	service := app.MigrationService{Writer: config.AtomicFileWriter{Path: filepath.Join(workdir, "weave.yaml")}}
	result, _, err := service.Migrate(ctx, cfg, dryRun)
	if err != nil {
		return exitCodeForError(app.WrapInvalidConfig(err)), app.WrapInvalidConfig(err)
	}

	if !result.Upgraded {
		fmt.Printf("migrate: no schema changes needed (version %d).\n", result.To)
		return ExitOK, nil
	}
	if result.DryRun {
		fmt.Printf("[dry-run] migrate: would upgrade weave.yaml schema from v%d to v%d.\n", result.From, result.To)
		return ExitOK, nil
	}
	fmt.Printf("migrate: upgraded weave.yaml schema from v%d to v%d.\n", result.From, result.To)
	return ExitOK, nil
}

func parseDryRunOnly(args []string, command string) (bool, error) {
	dryRun := false
	for _, a := range args {
		if a == "--dry-run" {
			dryRun = true
			continue
		}
		if strings.HasPrefix(a, "--") {
			return false, fmt.Errorf("unsupported flag for %s: %s", command, a)
		}
		return false, fmt.Errorf("unsupported argument for %s: %s", command, a)
	}
	return dryRun, nil
}

func printForgeSummary(result app.ForgeResult) {
	if result.DryRun {
		fmt.Printf("[dry-run] forge: planned %d filesystem change(s); no changes applied. Next step: rerun `weave forge` without --dry-run.\n", result.OpsPlanned)
		return
	}

	if result.WasNoOp {
		fmt.Println("forge: no changes needed; project is already converged.")
		return
	}

	fmt.Printf("forge: applied %d/%d filesystem change(s), weave.yaml updated.\n", result.OpsApplied, result.OpsPlanned)
}

func runProvider(ctx context.Context, args []string) (int, error) {
	action, name, dryRun, err := parseProviderAction(args)
	if err != nil {
		return exitCodeForError(err), err
	}

	workdir := resolveWorkdir()
	svc := newProviderService(workdir, nil)
	registry := providers.NewDefaultRegistry()

	result, names, err := runProviderAction(ctx, svc, registry, action, name, dryRun)
	if err != nil {
		return exitCodeForError(err), err
	}

	if action == providerActionList {
		for _, name := range names {
			fmt.Println(name)
		}
		return ExitOK, nil
	}

	printProviderSummary(action, result)

	return ExitOK, nil
}

func runAssetAdd(ctx context.Context, kind assetKind, name string, rest []string) (int, error) {
	fromFlag, dryRun, conflictFlag, err := parseAddFlags(rest)
	if err != nil {
		return exitCodeForError(err), err
	}

	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return exitCodeForError(app.WrapInvalidConfig(err)), app.WrapInvalidConfig(err)
	}

	service := newDefaultAssetAddService()
	policy := app.ConflictPolicy(cfg.Sync.ConflictPolicy)
	if conflictFlag != "" {
		policy = conflictFlag
	}

	result, err := service.Add(ctx, kind, name, fromFlag, dryRun, policy, cfg)
	if err != nil {
		return exitCodeForError(err), err
	}

	printAssetAddSummary(kind, name, result)

	return ExitOK, nil
}

func runDoctor(ctx context.Context, args []string) (int, error) {
	jsonOutput, err := parseJSONFlag(args)
	if err != nil {
		return exitCodeForError(err), err
	}

	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		wrapped := app.WrapInvalidConfig(err)
		return exitCodeForError(wrapped), wrapped
	}

	result, err := (app.DoctorService{ProviderRegistry: providers.NewDefaultRegistry()}).Run(ctx, workdir, cfg)
	if err != nil {
		return exitCodeForError(err), err
	}

	if jsonOutput {
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return exitCodeForError(err), err
		}
		fmt.Println(string(b))
	} else {
		printDoctorText(result)
	}

	if result.Status == app.DoctorStatusIssuesFound {
		return ExitDoctorIssues, nil
	}

	return ExitOK, nil
}

func parseJSONFlag(args []string) (bool, error) {
	for _, a := range args {
		if a == "--json" {
			return true, nil
		}
		if strings.HasPrefix(a, "--") {
			return false, fmt.Errorf("unsupported flag: %s", a)
		}
	}
	return false, nil
}

func printDoctorText(result app.DoctorResult) {
	if result.Status == app.DoctorStatusHealthy {
		fmt.Println("Doctor status: healthy")
		fmt.Println("No issues found. Your project is consistent with weave.yaml.")
		return
	}

	fmt.Println("Doctor status: issues_found")
	for idx, issue := range result.Issues {
		fmt.Printf("%d. [%s] %s\n", idx+1, issue.Code, issue.Summary)
		if issue.RepairCommand != "" {
			fmt.Printf("   Repair: %s\n", issue.RepairCommand)
		}
		fmt.Printf("   Docs: %s (%s)\n", issue.DocsPath, issue.DocsURL)
	}

	if len(result.RepairCommands) > 0 {
		fmt.Println("Suggested repair path:")
		for _, cmd := range result.RepairCommands {
			fmt.Printf("- %s\n", cmd)
		}
	}
}

func parseFromFlag(args []string) (string, error) {
	from, _, _, err := parseAddFlags(args)
	return from, err
}

func parseAddFlags(args []string) (string, bool, app.ConflictPolicy, error) {
	from := ""
	dryRun := false
	policy := app.ConflictPolicy("")

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--dry-run" {
			dryRun = true
			continue
		}

		if a == "--from" {
			if i+1 >= len(args) {
				return "", false, "", fmt.Errorf("--from requires a value")
			}
			from = args[i+1]
			i++
			continue
		}

		if strings.HasPrefix(a, "--from=") {
			from = strings.TrimPrefix(a, "--from=")
			continue
		}

		if a == "--overwrite" || a == "--skip" || a == "--backup" {
			next := app.ConflictPolicy(strings.TrimPrefix(a, "--"))
			if policy != "" && policy != next {
				return "", false, "", fmt.Errorf("only one conflict policy flag is allowed: --overwrite, --skip, or --backup")
			}
			policy = next
			continue
		}

		if strings.HasPrefix(a, "--") {
			return "", false, "", fmt.Errorf("unsupported flag: %s", a)
		}

		return "", false, "", fmt.Errorf("unsupported argument: %s", a)
	}

	return from, dryRun, policy, nil
}

func parseProviderAction(args []string) (providerAction, string, bool, error) {
	if len(args) == 0 {
		return "", "", false, fmt.Errorf("provider action is required: add|list|remove|repair")
	}

	dryRun := false
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--dry-run" {
			dryRun = true
			continue
		}
		rest = append(rest, a)
	}
	if len(rest) == 0 {
		return "", "", false, fmt.Errorf("provider action is required: add|list|remove|repair")
	}

	action := providerAction(rest[0])
	switch action {
	case providerActionList:
		if len(rest) > 1 {
			return "", "", false, fmt.Errorf("provider list does not accept a provider name")
		}
		if dryRun {
			return "", "", false, fmt.Errorf("--dry-run is only valid for mutating provider actions: add|remove|repair")
		}
		return action, "", false, nil
	case providerActionAdd, providerActionRemove, providerActionRepair:
		if len(rest) < 2 || strings.TrimSpace(rest[1]) == "" {
			return "", "", false, fmt.Errorf("provider name is required for %s", action)
		}
		if len(rest) > 2 {
			return "", "", false, fmt.Errorf("unsupported extra arguments for provider %s", action)
		}
		return action, rest[1], dryRun, nil
	default:
		return "", "", false, fmt.Errorf("unsupported provider action: %s", rest[0])
	}
}

func printAssetAddSummary(kind assetKind, name string, result app.AddAssetResult) {
	entity := string(kind)
	if result.DryRun {
		fmt.Printf("[dry-run] %s add %s: planned %d filesystem change(s); no changes applied. Next step: rerun without --dry-run.\n", entity, name, result.OpsPlanned)
		return
	}

	if result.ConflictDetected {
		if result.ConflictAction == string(app.ConflictPolicySkip) {
			fmt.Printf("%s add %s: conflict detected, skipped by policy. No changes applied.\n", entity, name)
			return
		}
		if result.ConflictAction == string(app.ConflictPolicyBackup) && result.BackupPath != "" {
			fmt.Printf("%s add %s: conflict detected, existing path backed up to %s before install.\n", entity, name, result.BackupPath)
		}
	}

	if result.ConfigSaved {
		fmt.Printf("%s add %s: applied %d/%d filesystem change(s), weave.yaml updated.\n", entity, name, result.OpsApplied, result.OpsPlanned)
		return
	}

	fmt.Printf("%s add %s: no config changes saved. Run `weave doctor` for repair guidance.\n", entity, name)
}

func printProviderSummary(action providerAction, result app.ProviderAddResult) {
	if result.DryRun {
		fmt.Printf("[dry-run] provider %s %s: planned %d filesystem change(s); no changes applied. Next step: rerun without --dry-run.\n", action, result.Provider, result.OpsPlanned)
		return
	}

	verb := string(action)
	if action == providerActionAdd {
		verb = "added"
	} else if action == providerActionRemove {
		verb = "removed"
	} else if action == providerActionRepair {
		verb = "repaired"
	}

	fmt.Printf("provider %s %s: applied %d/%d filesystem change(s), weave.yaml updated.\n", verb, result.Provider, result.OpsApplied, result.OpsPlanned)
}

func printHelp() {
	fmt.Println("Weave CLI — Declarative project control plane")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  weave [forge]")
	fmt.Println("  weave skill add <name> [--from <dir>] [--overwrite|--skip|--backup] [--dry-run]")
	fmt.Println("  weave command add <name> [--from <dir>] [--overwrite|--skip|--backup] [--dry-run]")
	fmt.Println("  weave provider <add|remove|repair|list> [name] [--dry-run]")
	fmt.Println("  weave doctor [--json]")
	fmt.Println("  weave migrate [--dry-run]")
	fmt.Println("  weave --version")
	fmt.Println("")
	fmt.Println("60-second quickstart:")
	fmt.Println("  1) weave forge")
	fmt.Println("  2) weave provider add claude-code")
	fmt.Println("  3) weave skill add sdd-orchestrator")
	fmt.Println("  4) weave doctor")
}

func printFirstRunGuidance() {
	fmt.Println("First run complete. Suggested next steps:")
	fmt.Println("- weave provider add claude-code")
	fmt.Println("- weave skill add sdd-orchestrator")
	fmt.Println("- weave doctor")
}

func configExists(root string) bool {
	if root == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(root, "weave.yaml"))
	return err == nil
}
