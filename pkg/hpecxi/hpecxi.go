package hpecxi

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

const HPEvendorID string = "0x17db"

func findLibs(libName, libPath string) ([]string, error) {
	var files []string

	fileInfos, err := os.ReadDir(libPath)
	if err != nil {
		klog.Errorf("Error while looking for %s in %s", libName, libPath)
		return nil, err
	}
	notFound := true
	for _, fileInfo := range fileInfos {
		if strings.HasPrefix(fileInfo.Name(), libName) {
			fullPath := filepath.Join(libPath, fileInfo.Name())
			files = append(files, fullPath)
			notFound = false
		}
	}
	if notFound {
		klog.Infof("Library %s not found at %s", libName, libPath)
	}
	return files, nil
}

// Support for initial function using hard coded defaults
// These can be further exposed if needed. Right now this is just for testing.
func GetDevices() map[string]int {

	// Note that these were the initial hard coded paths
	cfg := &HPECXIConfig{
		CxiDriverRoot:   "/sys/module/cxi_core/drivers",
		LibfabricPath:   "/opt/cray/lib64",
		LibcxiPath:      LibcxiPath,
		NetDevicePrefix: "hsn",
		PCIName:         "pci:cxi_core",
	}

	// Create a new plugin manager for the lister
	mgr := NewManager(cfg)
	return mgr.GetDevices()
}

// GetHPECXIs return a map of HPE Cassini on a node identified by the part of the pci address
// This may be changed to use cxilib calls instead of sysfs.
func (m *Manager) GetDevices() map[string]int {
	if _, err := os.Stat(m.config.CxiDriverRoot); err != nil {
		klog.Warningf("HPE CXI driver unavailable: %s", err)
		return make(map[string]int)
	}
	regex := m.GetCXIDriverRegex()
	matches, _ := filepath.Glob(regex)
	klog.Info(matches)
	devices := make(map[string]int)

	for _, path := range matches {
		klog.Info(path)
		devPaths, _ := filepath.Glob(path + "/net/*")

		for _, devPath := range devPaths {
			name := filepath.Base(devPath)
			if name[0:3] == m.config.NetDevicePrefix {
				nic_id, _ := strconv.Atoi(name[len(name)-1:])
				devices[name] = nic_id
			}
		}

	}

	for device, _ := range devices {
		klog.Info("Found device ", device)
	}

	return devices
}

// HPECXI check if a particular card is a HPE CXI NIC by checking the device's vendor ID
func HPECXI(cardName string) bool {
	sysfsVendorPath := "/sys/class/cxi/" + cardName + "/device/vendor"
	b, err := os.ReadFile(sysfsVendorPath)
	if err == nil {
		vid := strings.TrimSpace(string(b))

		if vid == HPEvendorID {
			return true
		} else {
			klog.Infof("%s is not a HPE NIC.", cardName)
		}
	} else {
		klog.Errorf("Error opening %s: %s", sysfsVendorPath, err)
	}
	return false
}
