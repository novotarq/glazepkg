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
	out, err := exec.Command("mise", "list").Output()
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
