package manager

import (
	"bufio"
	"os/exec"
	"strings"
	"time"

	"github.com/neur0map/glazepkg/internal/model"
)

type Mise struct{}

func (y *Mise) Name() model.Source { return model.SourceMise }

func (y *Mise) Available() bool { return commandExists("mise") }

func (y *Mise) Scan() ([]model.Package, error) {
	out, err := exec.Command("mise", "ls").Output()
	if err != nil {
		return nil, err
	}

	var pkgs []model.Package
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		pkgs = append(pkgs, model.Package{
			Name:        fields[0],
			Version:     fields[1],
			Source:      model.SourceMise,
			InstalledAt: time.Now(),
		})
	}
	return pkgs, nil
}

func (y *Mise) Search(query string) ([]model.Package, error) {
	out, err := exec.Command("mise", "search", query).Output()
	if err != nil || len(out) == 0 {
		return nil, nil
	}

	var pkgs []model.Package
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mise ERROR") {
			break
		}
		if strings.Contains(line, "Tool") && strings.Contains(line, "Description") {
			continue
		}
		tokens := strings.Split(line, "  ")
		name := strings.TrimSpace(tokens[0])
		description := strings.TrimSpace(tokens[len(tokens)-1])

		pkgs = append(pkgs, model.Package{Name: name, Description: description, Source: model.SourceMise})
	}

	return pkgs, nil
}
