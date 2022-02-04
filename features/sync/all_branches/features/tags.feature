Feature: sync all branches syncs the tags

  Scenario:
    Given my repo has the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    And I am on the "main" branch
    When I run "git-town sync --all"
    Then my repo now has the tags
      | NAME       | LOCATION      |
      | local-tag  | local, remote |
      | remote-tag | local, remote |
