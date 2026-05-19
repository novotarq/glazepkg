package manager

import (
	"bufio"
	"os/exec"
	"strings"

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
	// skip the description
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		name := strings.TrimSpace(fields[0])
		version := strings.TrimSpace(fields[1])

		pkgs = append(pkgs, model.Package{
			Name:    name,
			Version: version,
			Source:  model.SourceMise,
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

	// skip the first line - it is either an error or description
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		name := strings.TrimSpace(tokens[0])
		description := strings.TrimSpace(tokens[1])

		pkgs = append(pkgs, model.Package{
			Name:        name,
			Description: description,
			Source:      model.SourceMise,
		})
	}

	return pkgs, nil
}

func (y *Mise) CheckUpdates(pkgs []model.Package) map[string]string {
	updates := make(map[string]string)

	var package_names strings.Builder
	for i, p := range pkgs {
		if i > 0 {
			package_names.WriteString(" ")
		}
		package_names.WriteString(p.Name)
	}

	query := package_names.String()

	out, err := exec.Command("mise", "outdated", query).Output()
	if err != nil || len(out) == 0 {
		return updates
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	// always skip the first line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		name := strings.TrimSpace(fields[0])
		version := strings.TrimSpace(fields[3])
		updates[name] = version
	}

	return updates
}

func (y *Mise) InstallCmd(name string) *exec.Cmd {
	return exec.Command("mise", "use", name)
}

func (y *Mise) RemoveCmd(name string) *exec.Cmd {
	return exec.Command("mise", "unuse", name)
}

func (y *Mise) UpgradeCmd(name string) *exec.Cmd {
	return exec.Command("mise", "upgrade", name)
}
