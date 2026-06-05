package node

import (
	ipfs_config "github.com/ipfs/kubo/config"
	ipfs_repo "github.com/ipfs/kubo/repo"
)

var _ ipfs_repo.Repo = (*RepoMobile)(nil)

type RepoConfigPatch func(cfg *ipfs_config.Config) (err error)

type RepoMobile struct {
	ipfs_repo.Repo

	repoPath string
}

func (mr *RepoMobile) Path() string { return mr.repoPath }

func NewRepoMobile(path string, repo ipfs_repo.Repo) *RepoMobile {
	return &RepoMobile{
		Repo:     repo,
		repoPath: path,
	}
}

func (mr *RepoMobile) ApplyPatchs(patchs ...RepoConfigPatch) error {
	cfg, err := mr.Config()
	if err != nil {
		return err
	}

	if err := ChainIpfsConfigPatch(patchs...)(cfg); err != nil {
		return err
	}

	return mr.SetConfig(cfg)
}

func ChainIpfsConfigPatch(patchs ...RepoConfigPatch) RepoConfigPatch {
	return func(cfg *ipfs_config.Config) (err error) {
		for _, patch := range patchs {
			if patch == nil {
				continue // skip empty patch
			}

			if err = patch(cfg); err != nil {
				return
			}
		}
		return
	}
}
