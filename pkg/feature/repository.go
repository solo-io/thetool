package feature

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	// ReposFileName represents the default filename for repositories
	ReposFileName = "repositories.json"
)

type Repository struct {
	URL      string `json:"url"`
	Commit   string `json:"commit"`
	Manifest string `json:"manifest"`
}

type RepositoryStore interface {
	Init() error
	Add(Repository) error
	Remove(string) error
	List() ([]Repository, error)
}

type FileRepoStore struct {
	Filename string
}

func (r *FileRepoStore) Init() error {
	return r.save([]Repository{})
}
func (r *FileRepoStore) AddOrUpdate(repo Repository) (bool, error) {
	repos, err := r.List()
	if err != nil {
		return false, err
	}

	alreadyExists := false
	for i, r := range repos {
		if r.URL == repo.URL {
			alreadyExists = true
			repos[i].Commit = repo.Commit
			break
		}
	}
	if !alreadyExists {
		repos = append(repos, repo)
	}
	return alreadyExists, r.save(repos)
}

func (r *FileRepoStore) Remove(repoURL string) error {
	existing, err := r.List()
	if err != nil {
		return err
	}

	updated := []Repository{}
	for _, e := range existing {
		if e.URL != repoURL {
			updated = append(updated, e)
		}
	}

	if len(updated) == len(existing) {
		return fmt.Errorf("did not find repository with url %s to remove", repoURL)
	}

	return r.save(updated)
}

func (r *FileRepoStore) save(repos []Repository) error {
	b, err := json.MarshalIndent(repoFile{
		Date:        time.Now(),
		GeneratedBy: "thetool",
		Repos:       repos,
	}, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.Filename, b, 0644)
}

func (r *FileRepoStore) List() ([]Repository, error) {
	b, err := ioutil.ReadFile(r.Filename)
	if err != nil {
		return nil, err
	}
	repos := &repoFile{}
	err = json.Unmarshal(b, repos)
	if err != nil {
		return nil, err
	}
	return repos.Repos, nil
}

type repoFile struct {
	Date        time.Time    `json:"date"`
	GeneratedBy string       `json:"generatedBy"`
	Repos       []Repository `json:"repositories"`
}
