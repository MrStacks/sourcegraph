package main

import (
	"strings"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/batches/types"
	"github.com/sourcegraph/sourcegraph/internal/gqltestutil"
	"github.com/sourcegraph/sourcegraph/lib/batches/execution"
)

func testHelloWorldBatchChange() Test {
	batchSpec := `
name: e2e-test-batch-change
description: Add Hello World to READMEs

on:
  - repository: github.com/sourcegraph/automation-testing
    # This branch is never changing - hopefully.
    branch: executors-e2e

steps:
  - run: IFS=$'\n'; echo Hello World | tee -a $(find -name README.md)
    container: ubuntu:18.04

changesetTemplate:
  title: Hello World
  body: My first batch change!
  branch: hello-world # Push the commit to this branch.
  commit:
    message: Append Hello World to all README.md files
`

	expectedDiff := strings.Join([]string{
		"diff --git README.md README.md",
		"index 1914491..89e55af 100644",
		"--- README.md",
		"+++ README.md",
		"@@ -3,4 +3,4 @@ This repository is used to test opening and closing pull request with Automation",
		" ",
		" (c) Copyright Sourcegraph 2013-2020.",
		" (c) Copyright Sourcegraph 2013-2020.",
		"-(c) Copyright Sourcegraph 2013-2020.",
		"\\ No newline at end of file",
		"+(c) Copyright Sourcegraph 2013-2020.Hello World",
		"diff --git examples/README.md examples/README.md",
		"index 40452a9..a32cc2f 100644",
		"--- examples/README.md",
		"+++ examples/README.md",
		"@@ -5,4 +5,4 @@ This folder contains examples",
		" (This is a test message, ignore)",
		" ",
		" (c) Copyright Sourcegraph 2013-2020.",
		"-(c) Copyright Sourcegraph 2013-2020.",
		"\\ No newline at end of file",
		"+(c) Copyright Sourcegraph 2013-2020.Hello World",
		"diff --git examples/project3/README.md examples/project3/README.md",
		"index 272d9ea..f49f17d 100644",
		"--- examples/project3/README.md",
		"+++ examples/project3/README.md",
		"@@ -1,4 +1,4 @@",
		" # project3",
		" ",
		" (c) Copyright Sourcegraph 2013-2020.",
		"-(c) Copyright Sourcegraph 2013-2020.",
		"\\ No newline at end of file",
		"+(c) Copyright Sourcegraph 2013-2020.Hello World",
		"diff --git project1/README.md project1/README.md",
		"index 8a5f437..6284591 100644",
		"--- project1/README.md",
		"+++ project1/README.md",
		"@@ -3,4 +3,4 @@",
		" This is project 1.",
		" ",
		" (c) Copyright Sourcegraph 2013-2020.",
		"-(c) Copyright Sourcegraph 2013-2020.",
		"\\ No newline at end of file",
		"+(c) Copyright Sourcegraph 2013-2020.Hello World",
		"diff --git project2/README.md project2/README.md",
		"index b1e1cdd..9445efe 100644",
		"--- project2/README.md",
		"+++ project2/README.md",
		"@@ -1,3 +1,3 @@",
		" This is a starter template for [Learn Next.js](https://nextjs.org/learn).",
		" (c) Copyright Sourcegraph 2013-2020.",
		"-(c) Copyright Sourcegraph 2013-2020.",
		"\\ No newline at end of file",
		"+(c) Copyright Sourcegraph 2013-2020.Hello World",
		"",
	}, "\n")

	expectedState := gqltestutil.BatchSpecDeep{
		State: "COMPLETED",
		ChangesetSpecs: gqltestutil.BatchSpecChangesetSpecs{
			TotalCount: 1,
			Nodes: []gqltestutil.ChangesetSpec{
				{
					Type: "BRANCH",
					Description: gqltestutil.ChangesetSpecDescription{
						BaseRepository: gqltestutil.ChangesetRepository{Name: "github.com/sourcegraph/automation-testing"},
						BaseRef:        "executors-e2e",
						BaseRev:        "1c94aaf85d51e9d016b8ce4639b9f022d94c52e6",
						HeadRef:        "hello-world",
						Title:          "Hello World",
						Body:           "My first batch change!",
						Commits: []gqltestutil.ChangesetSpecCommit{
							{
								Message: "Append Hello World to all README.md files",
								Subject: "Append Hello World to all README.md files",
								Body:    "",
								Author: gqltestutil.ChangesetSpecCommitAuthor{
									Name:  "sourcegraph",
									Email: "sourcegraph@sourcegraph.com",
								},
							},
						},
						Diffs: gqltestutil.ChangesetSpecDiffs{
							FileDiffs: gqltestutil.ChangesetSpecFileDiffs{
								RawDiff: ``,
							},
						},
					},
					ForkTarget: gqltestutil.ChangesetForkTarget{},
				},
			},
		},
		Namespace: gqltestutil.Namespace{},
		WorkspaceResolution: gqltestutil.WorkspaceResolution{
			Workspaces: gqltestutil.WorkspaceResolutionWorkspaces{
				TotalCount: 1,
				Stats: gqltestutil.WorkspaceResolutionWorkspacesStats{
					Completed: 1,
				},
				Nodes: []gqltestutil.BatchSpecWorkspace{
					{
						State: "COMPLETED",
						DiffStat: gqltestutil.DiffStat{
							Added:   5,
							Deleted: 5,
						},
						Repository: gqltestutil.ChangesetRepository{Name: "github.com/sourcegraph/automation-testing"},
						Branch: gqltestutil.WorkspaceBranch{
							Name: "executors-e2e",
						},
						ChangesetSpecs: []gqltestutil.WorkspaceChangesetSpec{
							{},
						},
						SearchResultPaths: []string{},
						Executor: gqltestutil.Executor{
							QueueName: "batches",
							Active:    true,
						},
						Stages: gqltestutil.BatchSpecWorkspaceStages{
							Setup: []gqltestutil.ExecutionLogEntry{
								{
									Key:      "setup.git.init",
									ExitCode: 0,
								},
								{
									Key:      "setup.git.add-remote",
									ExitCode: 0,
								},
								{
									Key:      "setup.git.disable-gc",
									ExitCode: 0,
								},
								{
									Key:      "setup.git.fetch",
									ExitCode: 0,
								},
								{
									Key:      "setup.git.checkout",
									ExitCode: 0,
								},
								{
									Key:      "setup.git.set-remote",
									ExitCode: 0,
								},
								{
									Key:      "setup.fs.extras",
									Command:  []string{},
									ExitCode: 0,
								},
							},
							SrcExec: []gqltestutil.ExecutionLogEntry{
								{
									Key:      "step.docker.step.0.pre",
									ExitCode: 0,
								},
								{
									Key:      "step.docker.step.0.run",
									ExitCode: 0,
								},
								{
									Key:      "step.docker.step.0.post",
									ExitCode: 0,
								},
							},
							Teardown: []gqltestutil.ExecutionLogEntry{
								{
									Key:      "teardown.fs",
									Command:  []string{},
									ExitCode: 0,
								},
							},
						},
						Steps: []gqltestutil.BatchSpecWorkspaceStep{
							{
								Number:    1,
								Run:       "IFS=$'\\n'; echo Hello World | tee -a $(find -name README.md)",
								Container: "ubuntu:18.04",
								OutputLines: gqltestutil.WorkspaceOutputLines{
									Nodes: []string{
										"stderr: WARNING: The requested image's platform (linux/amd64) does not match the detected host platform (linux/arm64/v8) and no specific platform was requested",
										"stderr: Hello World",
										"",
									},
									TotalCount: 3,
								},
								ExitCode:        0,
								Environment:     []gqltestutil.WorkspaceEnvironmentVariable{},
								OutputVariables: []gqltestutil.WorkspaceOutputVariable{},
								DiffStat: gqltestutil.DiffStat{
									Added:   5,
									Deleted: 5,
								},
								Diff: gqltestutil.ChangesetSpecDiffs{
									FileDiffs: gqltestutil.ChangesetSpecFileDiffs{
										RawDiff: "diff --git README.md README.md\nindex 1914491..89e55af 100644\n--- README.md\n+++ README.md\n@@ -3,4 +3,4 @@ This repository is used to test opening and closing pull request with Automation\n \n (c) Copyright Sourcegraph 2013-2020.\n (c) Copyright Sourcegraph 2013-2020.\n-(c) Copyright Sourcegraph 2013-2020.\n\\ No newline at end of file\n+(c) Copyright Sourcegraph 2013-2020.Hello World\ndiff --git examples/README.md examples/README.md\nindex 40452a9..a32cc2f 100644\n--- examples/README.md\n+++ examples/README.md\n@@ -5,4 +5,4 @@ This folder contains examples\n (This is a test message, ignore)\n \n (c) Copyright Sourcegraph 2013-2020.\n-(c) Copyright Sourcegraph 2013-2020.\n\\ No newline at end of file\n+(c) Copyright Sourcegraph 2013-2020.Hello World\ndiff --git examples/project3/README.md examples/project3/README.md\nindex 272d9ea..f49f17d 100644\n--- examples/project3/README.md\n+++ examples/project3/README.md\n@@ -1,4 +1,4 @@\n # project3\n \n (c) Copyright Sourcegraph 2013-2020.\n-(c) Copyright Sourcegraph 2013-2020.\n\\ No newline at end of file\n+(c) Copyright Sourcegraph 2013-2020.Hello World\ndiff --git project1/README.md project1/README.md\nindex 8a5f437..6284591 100644\n--- project1/README.md\n+++ project1/README.md\n@@ -3,4 +3,4 @@\n This is project 1.\n \n (c) Copyright Sourcegraph 2013-2020.\n-(c) Copyright Sourcegraph 2013-2020.\n\\ No newline at end of file\n+(c) Copyright Sourcegraph 2013-2020.Hello World\ndiff --git project2/README.md project2/README.md\nindex b1e1cdd..9445efe 100644\n--- project2/README.md\n+++ project2/README.md\n@@ -1,3 +1,3 @@\n This is a starter template for [Learn Next.js](https://nextjs.org/learn).\n (c) Copyright Sourcegraph 2013-2020.\n-(c) Copyright Sourcegraph 2013-2020.\n\\ No newline at end of file\n+(c) Copyright Sourcegraph 2013-2020.Hello World\n",
									},
								},
							},
						},
					},
				},
			},
		},
		Source: "REMOTE",
		Files: gqltestutil.BatchSpecFiles{
			TotalCount: 0,
			Nodes:      []gqltestutil.BatchSpecFile{},
		},
	}

	return Test{
		PreExistingCacheEntries: map[string]execution.AfterStepResult{},
		BatchSpecInput:          batchSpec,
		ExpectedCacheEntries: map[string]execution.AfterStepResult{
			"wHcoEItqNkdJj9k1-1sRCQ-step-0": {
				Version: 2,
				Stdout:  "Hello World\n",
				Diff:    []byte(expectedDiff),
				Outputs: map[string]any{},
			},
		},
		ExpectedChangesetSpecs: []*types.ChangesetSpec{
			{
				Type:              "branch",
				DiffStatAdded:     5,
				DiffStatDeleted:   5,
				BatchSpecID:       2,
				BaseRepoID:        1,
				UserID:            1,
				BaseRev:           "1c94aaf85d51e9d016b8ce4639b9f022d94c52e6",
				BaseRef:           "executors-e2e",
				HeadRef:           "refs/heads/hello-world",
				Title:             "Hello World",
				Body:              "My first batch change!",
				CommitMessage:     "Append Hello World to all README.md files",
				CommitAuthorName:  "sourcegraph",
				CommitAuthorEmail: "sourcegraph@sourcegraph.com",
				Diff:              []byte(expectedDiff),
			},
		},
		ExpectedState: expectedState,
		CacheDisabled: true,
	}
}