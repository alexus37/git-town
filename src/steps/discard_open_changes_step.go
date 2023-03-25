package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	EmptyStep
}

func (step *DiscardOpenChangesStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.DiscardOpenChanges()
}
