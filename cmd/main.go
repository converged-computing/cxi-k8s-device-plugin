package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/HewlettPackard/cxi-k8s-device-plugin/pkg/hpecxi"
	"github.com/HewlettPackard/cxi-k8s-device-plugin/pkg/plugin"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"k8s.io/klog/v2"
)

var version string

func main() {

	versions := [...]string{
		"HPE Slingshot device plugin for Kubernetes",
		fmt.Sprintf("%s version %s", os.Args[0], version),
	}

	flag.Usage = func() {
		for _, v := range versions {
			fmt.Fprintf(os.Stderr, "%s\n", v)
		}
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
	}
	var pulse int
	var libfabricPath, libcxiPath, pciName, netDevicePrefix, cxiDriverRoot string
	flag.IntVar(&pulse, "pulse", 0, "time between health check polling in seconds.  Set to 0 to disable.")
	flag.StringVar(&netDevicePrefix, "net-device", hpecxi.NetDevicePrefix, "Device prefix to search for in net (e.g, hsi)")
	flag.StringVar(&libfabricPath, "libfabric", hpecxi.LibfabricPath, "Directory path to lib64 with libfabric")
	flag.StringVar(&libcxiPath, "libcxi", hpecxi.LibcxiPath, "Directory path to lib64 with libfabric")
	flag.StringVar(&pciName, "pci-name", hpecxi.PCIName, "PCI device name (e.g, pci:cxi_ss1")
	flag.StringVar(&cxiDriverRoot, "cxi-driver-root", hpecxi.CxiDriverRoot, "/sys/modules/<x>/devices root")
	flag.Parse()

	for _, v := range versions {
		klog.Infof("%s", v)
	}

	// Configuration for paths, naming
	cfg := &hpecxi.HPECXIConfig{
		CxiDriverRoot:   cxiDriverRoot,
		LibfabricPath:   libfabricPath,
		LibcxiPath:      libcxiPath,
		NetDevicePrefix: netDevicePrefix,
		PCIName:         pciName,
	}

	// Create a new plugin manager for the lister
	mgr := hpecxi.NewManager(cfg)
	l := plugin.NewHPECXILister(mgr)
	manager := dpm.NewManager(&l)

	// Tell user the configuration found
	klog.Info("🌊 Configuration:")
	klog.Infof("    Net Device Prefix: %s\n", cfg.NetDevicePrefix)
	klog.Infof("    CXI Driver Root:   %s\n", cfg.CxiDriverRoot)
	klog.Infof("    Libfabric Path:    %s\n", cfg.LibfabricPath)
	klog.Infof("    Libcxi Path:       %s\n", cfg.LibcxiPath)
	klog.Infof("    PCI Name:          %s\n", cfg.PCIName)

	if pulse > 0 {
		go func() {
			klog.Infof("Heart beating every %d seconds", pulse)
			for {
				time.Sleep(time.Second * time.Duration(pulse))
				l.Heartbeat <- true
			}
		}()
	}

	go func() {
		// Check if there are Cassini NICs installed
		// TODO: check how many and update channel.
		var path = "/sys/class/cxi"
		if _, err := os.Stat(path); err == nil {
			l.ResUpdateChan <- []string{"cxi"}
		}
	}()
	manager.Run()

}
