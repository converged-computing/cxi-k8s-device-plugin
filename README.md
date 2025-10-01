# CXI Kubernetes device plugin
![Tests](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/caio-davi/40551c54cd6c4079844a78fd6b46211f/raw/go-tests-badge.json&cacheSeconds=300)
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/caio-davi/40551c54cd6c4079844a78fd6b46211f/raw/coverage-badge.json&cacheSeconds=300)
[![Go Report Card](https://goreportcard.com/badge/github.com/HewlettPackard/cxi-k8s-device-plugin)](https://goreportcard.com/report/github.com/HewlettPackard/cxi-k8s-device-plugin)

> [!CAUTION]
This is a beta software, not recommended for production systems.

## Introduction 

This is a device plugin implementation that enables the registration of HPE Slingshot Cassini NICs in a Kubernetes cluster for compute workload.

## Requirements

##### - Kubernetes v1.30 or higher

##### - CXI Services Enabled

This step must be executed for every node and every Slingshot NIC that serves the Kubernetes cluster:
```
cxi_service enable -d cxiX -s 1
```

##### - Enable Mutating Admission Policy

Mutating Admission Policy is an alpha feature released in Kubernetes v1.30. As with any alpha feature, it must be explicitly enabled before usage. Add the following lines to /etc/kubernetes/manifests/kube-apiserver.yaml:
```
- --feature-gates=MutatingAdmissionPolicy=true
- --runtime-config=admissionregistration.k8s.io/v1alpha1=true
```

##### - Multus CNI 

Since we need multiple NICs, we are going to use Multus CNI (a CNI meta-plugin) that enables attaching multiple network interfaces to pods.
```
kubectl apply -f https://raw.githubusercontent.com/k8snetworkplumbingwg/multus-cni/master/deployments/multus-daemonset-thick.yml
```

##### - Whereabouts CNI

Whereabouts is an IP Address Management (IPAM) CNI plugin that assigns IP addresses cluster-wide.  
```
git clone https://github.com/k8snetworkplumbingwg/whereabouts
cd whereabouts
kubectl apply \
    -f doc/crds/daemonset-install.yaml \
    -f doc/crds/whereabouts.cni.cncf.io_ippools.yaml \
    -f doc/crds/whereabouts.cni.cncf.io_overlappingrangeipreservations.yaml
```

## Deployment

Alpha-version image available at `hub.docker.hpecorp.net/caio.davi/cxi-k8s-device-plugin:0.1`. 
The device plugin should run in every node with available Cassini NICs, in order to make them available to the K8s cluster. The easiest way of doing so is to create a Kubernetes DaemonSet:

```
kubectl apply \
    -f ./deploy/NetworkAttachmentDefinition \
    -f ./deploy/MutatingAdmissionPolicy \ 
    -f ./deploy/hpecxi-device-plugin-ds.yaml
```

## Development

You can build the docker container:

```bash
DOCKER_TAG=04 make docker-build
DOCKER_TAG=04 make docker-push
```

## Example

### Pod Resources

Here is an example for how to request a device from a PodSpec

```yaml
resources:
  requests: 
    beta.hpe.com/cxi: 1
```

### Running

You can run an example with a Flux MiniCluster.

```bash
# Look at this file first an ensure you have updated the entrypoint for your setup
kubectl apply -f ./deploy/hpecxi-device-plugin-ds.yaml 
```

Here are the arguments you can set:

```console
/go/bin/cxi-k8s-device-plugin --help
HPE Slingshot device plugin for Kubernetes
/go/bin/cxi-k8s-device-plugin version 0.0.1-beta
Usage:
  -alsologtostderr
    	log to standard error as well as files
  -cxi-driver-root string
    	/sys/modules/<x>/devices root (default "/sys/module/cxi_ss1/drivers")
  -libcxi string
    	Directory path to lib64 with libfabric (default "/usr/lib64")
  -libfabric string
    	Directory path to lib64 with libfabric (default "/opt/cray/libfabric/2.1/lib64")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -log_link string
    	If non-empty, add symbolic links in this directory to the log files
  -logbuflevel int
    	Buffer log messages logged at this level or lower (-1 means don't buffer; 0 means buffer INFO only; ...). Has limited applicability on non-prod platforms.
  -logtostderr
    	log to standard error instead of files
  -net-device string
    	Device prefix to search for in net (e.g, hsi) (default "hsi")
  -pci-name string
    	PCI device name (e.g, pci:cxi_ss1 (default "pci:cxi_ss1")
  -pulse int
    	time between health check polling in seconds.  Set to 0 to disable.
  -stderrthreshold value
    	logs at or above this threshold go to stderr (default 2)
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

Create the MiniCluster (non interactive)

```bash
kubectl apply -f example/flux-minicluster.yaml
```

And watch lammps run!

<details>

<summary>LAMMPS Log</summary>

```console
[sochat1@hetchy1001:deploy]$ kubectl logs lmp-0-8tjgt -f
Defaulted container "lmp" out of: lmp, flux-view (init)
🟧️  wait-fs: 2025/10/01 07:15:11 wait-fs.go:40: /mnt/flux/flux-operator-done.txt
🟧️  wait-fs: 2025/10/01 07:15:11 wait-fs.go:49: Found existing path /mnt/flux/flux-operator-done.txt

Hello user root

🌟️ Curve Certificate
curve.cert
#   ****  Generated on 2023-04-26 22:54:42 by CZMQ  ****
#   ZeroMQ CURVE **Secret** Certificate
#   DO NOT PROVIDE THIS FILE TO OTHER USERS nor change its permissions.
    
metadata
    name = "flux-cert-generator"
    keygen.hostname = "lmp-0"
curve
    public-key = "5*NS#QbaV-ean:38}mN+I1FrcetR9cuFRLDhC?Hf"
    secret-key = "goN&y=}!Vn(nt7G4Zo-MCpiU[TwYW&3#X&t<:!qJ"

📦 Resources
flux R encode --hosts=lmp-[0-1] --local
{"version": 1, "execution": {"R_lite": [{"rank": "0-1", "children": {"core": "0-63"}}], "starttime": 0.0, "expiration": 0.0, "nodelist": ["lmp-[0-1]"]}}
👋 Hello, I'm lmp-0
The main host is lmp-0
The working directory is /opt/lammps/examples/reaxff/HNS, contents include:
README.txt	ffield.reax.hns  log.30Nov23.reaxff.hns.g++.1
data.hns-equil	in.reaxff.hns	 log.30Nov23.reaxff.hns.g++.4
🚩️ Flux Option Flags defined
Command provided is: lmp -v x 8 -v y 8 -v z 8 -in in.reaxff.hns -nocite
Flags for flux are -N 2 -n128  

🌀 Submit Mode: flux start -o --config /mnt/flux/view/etc/flux/config -Scron.directory=/etc/flux/system/cron.d   -Stbon.fanout=256   -Srundir=/mnt/flux/view/run/flux    -Sstatedir=/mnt/flux/view/var/lib/flux -Slocal-uri=local:///mnt/flux/view/run/flux/local -Stbon.connect_timeout=5s     -Slog-stderr-level=6    -Slog-stderr-mode=local  flux submit  -N 2 -n128   --quiet --watch lmp -v x 8 -v y 8 -v z 8 -in in.reaxff.hns -nocite
Flags for flux are -N 2 -n128  
broker.info[0]: start: none->join 0.433807ms
broker.info[0]: parent-none: join->init 0.017554ms
cron.info[0]: synchronizing cron tasks to event heartbeat.pulse
job-manager.info[0]: restart: 0 jobs
job-manager.info[0]: restart: 0 running jobs
job-manager.info[0]: restart: checkpoint.job-manager not found
broker.info[0]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
sched-fluxion-resource.info[0]: version 0.45.0
sched-fluxion-resource.warning[0]: create_reader: allowlist unsupported
sched-fluxion-resource.info[0]: populate_resource_db: loaded resources from core's resource.acquire
sched-fluxion-qmanager.info[0]: version 0.45.0
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: tab: cron-1 created: scheduled in 71088.204s at Thu Oct  2 03:00:00 2025
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 0.4s
broker.info[0]: rc1-success: init->quorum 0.397933s
broker.info[0]: online: lmp-0 (ranks 0)
broker.info[0]: online: lmp-[0-1] (ranks 0-1)
broker.info[0]: quorum-full: quorum->run 0.448331s
LAMMPS (22 Jul 2025 - Development - patch_22Jul2025-382-g1db2e93763)
OMP_NUM_THREADS environment is not set. Defaulting to 1 thread.
  using 1 OpenMP thread(s) per MPI task
Reading data file ...
  triclinic box = (0 0 0) to (22.326 11.1412 13.778966) with tilt (0 -5.02603 0)
  8 by 4 by 4 MPI processor grid
  reading atoms ...
  304 atoms
  reading velocities ...
  304 velocities
  read_data CPU = 0.052 seconds
Replication is creating a 8x8x8 = 512 times larger system...
  triclinic box = (0 0 0) to (178.608 89.1296 110.23173) with tilt (0 -40.20824 0)
  8 by 4 by 4 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 19.79 out of 512 (3.87%)
  155648 atoms
  replicate CPU = 0.003 seconds
Neighbor list info ...
  update: every = 20 steps, delay = 0 steps, check = no
  max neighbors/atom: 2000, page size: 100000
  master list distance cutoff = 11
  ghost atom cutoff = 11
  binsize = 5.5, bins = 40 17 21
  2 neighbor lists, perpetual/occasional/extra = 2 0 0
  (1) pair reaxff, perpetual
      attributes: half, newton off, ghost
      pair build: half/bin/ghost/newtoff
      stencil: full/ghost/bin/3d
      bin: standard
  (2) fix qeq/reax, perpetual, copy from (1)
      attributes: half, newton off
      pair build: copy
      stencil: none
      bin: none
Setting up Verlet run ...
  Unit style    : real
  Current step  : 0
  Time step     : 0.1
Per MPI rank memory allocation (min/avg/max) = 143.9 | 143.9 | 143.9 Mbytes
   Step          Temp          PotEng         Press          E_vdwl         E_coul         Volume    
         0   300           -113.27833      438.99618     -111.57687     -1.7014647      1754807.5    
        10   300.64265     -113.28007      771.21336     -111.57866     -1.7014067      1754807.5    
        20   302.23163     -113.28471      1617.9776     -111.58344     -1.7012699      1754807.5    
        30   302.52602     -113.28543      4311.9345     -111.58441     -1.701021       1754807.5    
        40   301.00893     -113.28084      6495.276      -111.58016     -1.7006791      1754807.5    
        50   298.22387     -113.27248      6671.9892     -111.57218     -1.7003023      1754807.5    
        60   295.54892     -113.26445      6412.5588     -111.56453     -1.699926       1754807.5    
        70   294.96528     -113.26266      7033.6801     -111.56311     -1.6995494      1754807.5    
        80   297.40591     -113.26991      8436.3516     -111.57073     -1.6991764      1754807.5    
        90   301.11971     -113.28098      9412.0446     -111.58214     -1.6988469      1754807.5    
       100   302.41516     -113.28478      10326.738     -111.58617     -1.6986109      1754807.5    
Loop time of 21.7324 on 128 procs for 100 steps with 155648 atoms

Performance: 0.040 ns/day, 603.677 hours/ns, 4.601 timesteps/s, 716.204 katom-step/s
98.4% CPU use with 128 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 10.754     | 12.247     | 13.543     |  15.6 | 56.36
Neigh   | 0.21365    | 0.21834    | 0.2225     |   0.7 |  1.00
Comm    | 0.69795    | 2.0177     | 3.5951     |  39.7 |  9.28
Output  | 0.0037101  | 0.026436   | 0.043183   |  10.4 |  0.12
Modify  | 7.1327     | 7.2212     | 7.3356     |   2.3 | 33.23
Other   |            | 0.001477   |            |       |  0.01

Nlocal:           1216 ave        1222 max        1207 min
Histogram: 1 0 4 5 27 14 39 18 16 4
Nghost:        7592.56 ave        7610 max        7578 min
Histogram: 3 9 13 23 31 19 16 10 3 1
Neighs:         432968 ave      434953 max      429964 min
Histogram: 1 0 6 8 21 26 25 24 13 4

Total # of neighbors = 55419905
Ave neighs/atom = 356.05922
Neighbor list builds = 5
Dangerous builds not checked
Total wall time: 0:00:22
broker.info[0]: rc2.0: flux submit -N 2 -n128 --quiet --watch lmp -v x 8 -v y 8 -v z 8 -in in.reaxff.hns -nocite Exited (rc=0) 27.2s
broker.info[0]: rc2-success: run->cleanup 27.2071s
broker.info[0]: cleanup.0: flux queue stop --quiet --all --nocheckpoint Exited (rc=0) 0.1s
broker.info[0]: cleanup.1: flux resource acquire-mute Exited (rc=0) 0.1s
broker.info[0]: cleanup.2: flux cancel --user=all --quiet --states RUN Exited (rc=0) 0.1s
broker.info[0]: cleanup.3: flux queue idle --quiet Exited (rc=0) 0.1s
broker.info[0]: cleanup-success: cleanup->shutdown 0.451006s
broker.info[0]: children-complete: shutdown->finalize 57.1783ms
broker.info[0]: rc3.0: running /etc/flux/rc3.d/01-sched-fluxion
broker.info[0]: rc3.0: /etc/flux/rc3 Exited (rc=0) 0.1s
broker.info[0]: rc3-success: finalize->goodbye 96.4793ms
broker.info[0]: goodbye: goodbye->exit 0.04662ms
```

</details>

You can set `logging.quiet: true` to only see the LAMMPS logs. See device plugin logs if needed. There is one pod deployed per node.

```bash
$ kubectl logs -n kube-system hpecxi-device-plugin-daemonset-gc2lx -f
```

<details>

<summary>Daemonset Pod Log (Plugin Device)</summary>

```console
I1001 07:02:08.992426       1 main.go:43] HPE Slingshot device plugin for Kubernetes
I1001 07:02:08.992543       1 main.go:43] /go/bin/cxi-k8s-device-plugin version 0.0.1-beta
I1001 07:02:08.992549       1 main.go:61] 🌊 Configuration:
I1001 07:02:08.992553       1 main.go:62]     Net Device Prefix: hsi
I1001 07:02:08.992558       1 main.go:63]     CXI Driver Root:   /sys/module/cxi_ss1/drivers
I1001 07:02:08.992562       1 main.go:64]     Libfabric Path:    /opt/cray/libfabric/2.1/lib64
I1001 07:02:08.992566       1 main.go:65]     Libcxi Path:       /usr/lib64
I1001 07:02:08.992570       1 main.go:66]     PCI Name:          pci:cxi_ss1
I1001 07:02:08.992576       1 manager.go:42] Starting device plugin manager
I1001 07:02:08.992607       1 manager.go:46] Registering for system signal notifications
I1001 07:02:08.992776       1 manager.go:52] Registering for notifications of filesystem changes in device plugin directory
I1001 07:02:08.992839       1 manager.go:60] Starting Discovery on new plugins
I1001 07:02:08.992861       1 manager.go:66] Handling incoming signals
I1001 07:02:08.992871       1 manager.go:71] Received new list of plugins: [cxi]
I1001 07:02:08.992897       1 manager.go:110] Adding a new plugin "cxi"
I1001 07:02:08.992915       1 plugin.go:64] cxi: Starting plugin server
I1001 07:02:08.992927       1 plugin.go:94] cxi: Starting the DPI gRPC server
I1001 07:02:08.993092       1 plugin.go:112] cxi: Serving requests...
I1001 07:02:08.993100       1 plugin.go:128] cxi: Registering the DPI with Kubelet
I1001 07:02:08.993422       1 plugin.go:140] cxi: Registration for endpoint beta.hpe.com_cxi
I1001 07:02:08.995004       1 hpecxi.go:63] [/sys/module/cxi_ss1/drivers/pci:cxi_ss1/0000:01:00.0 /sys/module/cxi_ss1/drivers/pci:cxi_ss1/0000:c2:00.0]
I1001 07:02:08.995019       1 hpecxi.go:67] /sys/module/cxi_ss1/drivers/pci:cxi_ss1/0000:01:00.0
I1001 07:02:08.995061       1 hpecxi.go:67] /sys/module/cxi_ss1/drivers/pci:cxi_ss1/0000:c2:00.0
I1001 07:02:08.995099       1 hpecxi.go:81] Found device hsi1
I1001 07:02:08.995103       1 hpecxi.go:81] Found device hsi0
I1001 07:02:08.995106       1 plugin.go:124] Found 2 HPE Slingshot NICs
I1001 07:02:37.506990       1 plugin.go:185] Mounting /usr/lib64/libcxi.a
I1001 07:02:37.507023       1 plugin.go:185] Mounting /usr/lib64/libcxi.la
I1001 07:02:37.507027       1 plugin.go:185] Mounting /usr/lib64/libcxi.so
I1001 07:02:37.507030       1 plugin.go:185] Mounting /usr/lib64/libcxi.so.1
I1001 07:02:37.507032       1 plugin.go:185] Mounting /usr/lib64/libcxi.so.1.5.0
I1001 07:02:37.507034       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.a
I1001 07:02:37.507037       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.la
I1001 07:02:37.507039       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so
I1001 07:02:37.507041       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so.0
I1001 07:02:37.507044       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so.0.0.0
I1001 07:02:37.507046       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.a
I1001 07:02:37.507048       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.la
I1001 07:02:37.507050       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so
I1001 07:02:37.507052       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so.0
I1001 07:02:37.507055       1 plugin.go:185] Mounting /usr/lib64/libcxiutils.so.0.0.0
I1001 07:02:37.507057       1 plugin.go:185] Mounting /opt/cray/libfabric/2.1/lib64/libfabric
I1001 07:02:37.507059       1 plugin.go:185] Mounting /opt/cray/libfabric/2.1/lib64/libfabric.a
I1001 07:02:37.507062       1 plugin.go:185] Mounting /opt/cray/libfabric/2.1/lib64/libfabric.so
I1001 07:02:37.507064       1 plugin.go:185] Mounting /opt/cray/libfabric/2.1/lib64/libfabric.so.1
I1001 07:02:37.507066       1 plugin.go:185] Mounting /opt/cray/libfabric/2.1/lib64/libfabric.so.1.18.2
I1001 07:02:37.507070       1 plugin.go:197] Allocating cxi1
```

</details>

Don't forget to clean up!

> #### Make sure the IPAM definitions in the `./deploy/NetworkAttachmentDefinition` are follwoing your cluster network requirements. 
