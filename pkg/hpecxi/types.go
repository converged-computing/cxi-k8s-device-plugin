package hpecxi

import "fmt"

// Defaults
var (
	CxiDriverRoot = "/sys/module/cxi_ss1/drivers"
	LibfabricPath = "/opt/cray/libfabric/2.1/lib64"
	LibcxiPath    = "/usr/lib64"

	NetDevicePrefix = "hsi"

	// This can take on a multitude of names, e.g., pci:cxi_core
	PCIName = "pci:cxi_ss1"
)

// Definition of a group of CXI devices
type Manager struct {
	config  *HPECXIConfig
	Devices map[string]int
}

// NewHPECXI returns a new HPECXI, a set of CXI devices and configuation
func NewManager(cfg *HPECXIConfig) *Manager {
	return &Manager{
		config: cfg,
	}
}

func (h *Manager) EnvVars() map[string]string {
	return map[string]string{
		"LD_LIBRARY_PATH": fmt.Sprintf("%s:%s", h.config.LibfabricPath, h.config.LibcxiPath),
	}
}

// GetLibs returns libraries discovered in LibPaths
func (h *Manager) GetLibs() ([]string, error) {
	var libs []string

	for libname, libpath := range h.LibPaths() {
		newLibs, err := findLibs(libname, libpath)
		if err != nil {
			return nil, err
		}
		libs = append(libs, newLibs...)
	}

	return libs, nil
}

// GetCXIDriver regex returns a regular expression to discovery devices
func (m *Manager) GetCXIDriverRegex() string {
	return fmt.Sprintf("%s/%s/[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]:*", m.config.CxiDriverRoot, m.config.PCIName)
}

func (m *Manager) LibPaths() map[string]string {
	return map[string]string{
		"libfabric":   m.config.LibfabricPath,
		"libcxi":      m.config.LibcxiPath,
		"libcxiutils": m.config.LibcxiPath,
	}
}

// HPECXIConfig holds configuration defaults
type HPECXIConfig struct {
	CxiDriverRoot   string
	LibfabricPath   string
	LibcxiPath      string
	NetDevicePrefix string
	PCIName         string
}
