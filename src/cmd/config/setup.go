package config

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const setupConfigDesc = "Prompts to setup your Git Town configuration"

func setupCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigDesc,
		Long:  cmdhelpers.Long(setupConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigSetup(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigSetup(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := loadSetupConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	newMainBranch, aborted, err := dialog.EnterMainBranch(config.localBranches.Names(), repo.Runner.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetMainBranch(newMainBranch)
	if err != nil {
		return err
	}
	newPerennialBranches, aborted, err := dialog.EnterPerennialBranches(config.localBranches.Names(), repo.Runner.PerennialBranches, repo.Runner.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	if slices.Compare(repo.Runner.PerennialBranches, newPerennialBranches) != 0 {
		return repo.Runner.SetPerennialBranches(newPerennialBranches)
	}
	return nil
}

type setupConfig struct {
	localBranches gitdomain.BranchInfos
	dialogInputs  dialog.TestInputs
}

func loadSetupConfig(repo *execute.OpenRepoResult, verbose bool) (setupConfig, bool, error) {
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
	})
	dialogInputs := dialog.LoadTestInputs(os.Environ())
	return setupConfig{
		localBranches: branchesSnapshot.Branches,
		dialogInputs:  dialogInputs,
	}, exit, err
}