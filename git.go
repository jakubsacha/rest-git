package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

type branch struct {
	Name string
	Sha  string
}

func checkoutRepos(config Config, repos map[string]*git.Repository) error {
	for _, repo := range config.Repositories {
		var err error
		path := fmt.Sprintf("./repos/%s", repo.Name)
		log.Println("Opening repository", repo.Url, "to:", path)
		repos[repo.Name], err = git.PlainOpen(path)
		if err != nil {
			log.Println("Warning:", err)
			log.Println("Checking out", repo.Name, "to:", repo.Url)
			err := os.RemoveAll(path)
			if err != nil {
				log.Println("Couldnt clean up before checkout:", err)
			}

			repos[repo.Name], err = git.PlainClone(path, true, &git.CloneOptions{
				URL:        repo.Url,
				NoCheckout: true,
			})
			if err != nil {
				return fmt.Errorf("couldn't clone %s: %v", repo.Name, err)
			}
		}
	}
	return nil
}

func fetchRepos(repos map[string]*git.Repository) []string {
	var details []string
	log.Println("Fetching changes from git")
	for name, repo := range repos {
		log.Println("Fetching changes for repo", name)
		err := repo.Fetch(&git.FetchOptions{})
		if err != nil {
			if reflect.TypeOf(err).String() == "NoErrAlreadyUpToDate" {
				details = append(details, "Repository %s status: Nothing to fetch")
			} else {
				details = append(details, fmt.Sprintf("Repository %s status: %v", name, err))
				log.Println(err)
			}
		} else {
			details = append(details, "Repository %s status: Updated")
		}
	}

	return details
}

func listRemoteRefs(repos map[string]*git.Repository, name, prefix string) ([]*branch, error) {
	var branches []*branch

	remote, err := repos[name].Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("remote not found: %v", err)
	}
	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("couldn't list remote refs: %v", err)
	}
	refPrefix := fmt.Sprintf("refs/%s/", prefix)
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]
		branches = append(branches, &branch{Name: branchName, Sha: ref.Hash().String()})
	}
	sort.Slice(branches, func(i, j int) bool { return branches[i].Name < branches[j].Name })
	return branches, nil
}
