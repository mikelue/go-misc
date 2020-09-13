package env

import (
	"fmt"
	"testing"

	"github.com/mikelue/go-misc/utils/runtime"
	"github.com/mikelue/go-misc/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var currentSrcDir string

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Env Suite")
}

func init() {
	currentSrcDir = runtime.CallerUtils.GetDirOfSource()
}

func newXdgSetup(xdgConfig, subdir string) []utils.RollbackContainer {
	finalDir := fmt.Sprintf("%s/%s", xdgConfig, subdir)

	return []utils.RollbackContainer {
		utils.RollbackContainerBuilder.NewDir(finalDir),
		utils.RollbackContainerBuilder.NewEnv(map[string]string {
			"XDG_CONFIG_HOME": xdgConfig,
		}),
	}
}
