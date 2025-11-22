# Nexus Milestone 1: Core Runtime Implementation

This milestone implements the low-level container runtime interface using `libcontainer`. We are bypassing high-level abstractions to talk directly to the kernel.


### The Scope

I have achieved process isolation. I can spawn a node, assign it a discrete chunk of RAM (in bytes) and CPU shares, and isolate its PID namespace. The container successfully executes an Alpine-based rootfs.

### Environment Requirements

Don't try to run this on Windows or macOS. It won't work. The syscalls aren't there.

* **OS:** Linux (Ubuntu 22.04+ recommended).
* **Kernel:** 5.x+ (Must support Cgroup v2).
* **Privileges:** Root (We are manipulating `/sys/fs/cgroup` and namespaces. `sudo` is mandatory).
* **Init System:** Systemd (We fight with it for cgroup control, but the code handles the `cgroup.procs` write path manually).

### Dependency Setup

If you don't have Go installed, fix that first.

```bash
# Remove old go
sudo rm -rf /usr/local/go

# Install a version that actually works (1.21+)
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Path setup
export PATH=$PATH:/usr/local/go/bin
```

### Test Environment Preparation

The runtime requires a valid root filesystem to pivot into. We use a minimal Alpine Linux rootfs. If this directory is empty, the `setns` syscall will fail with `exit status 1`.

**Do not skip this.**

```bash
# 1. Create the storage location
sudo mkdir -p /var/lib/nexus/images/alpine-base

# 2. Fetch the mini rootfs
wget https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.4-x86_64.tar.gz

# 3. Extract. Note the -C. If you mess up the path, the container dies.
sudo tar -xvf alpine-minirootfs-3.18.4-x86_64.tar.gz -C /var/lib/nexus/images/alpine-base
```

### Build and Execution

We aren't using `libcontainer.New()` for the factory because it's deprecated logic. We use `libcontainer.Create()` and handle the Cgroup hierarchy creation manually to bypass Systemd's user slice restrictions.

```bash
# 1. Build the binary
go build -o nexus

# 2. Run the creation command (ROOT REQUIRED)
# We allocate 256MB RAM and 512 CPU shares.
sudo ./nexus node create node-1 --mem 256 --cpu 512
```

### Verification & Proof

The following screenshot demonstrates:

1.  Successful creation of the Cgroup hierarchy at `/sys/fs/cgroup/nexus`.
2.  Controller activation (`+cpu +memory +pids`) via manual write to `subtree_control`.
3.  The process launch returning a valid host PID.
4.  Verification via `ps` that the `sleep 3600` payload is running.

### Implementation Notes

* **Cgroup V2:** We are forcing the creation of `/sys/fs/cgroup/nexus` and explicitly enabling controllers. If your kernel is configured without `cgroup_no_v1=all` or similar hybrid nonsense, you might see issues. We assume pure V2.
* **PID -1 Fix:** We fetch the PID directly from the `process` struct post-launch, rather than relying on `container.State()`, which has a race condition on read.
* **Isolation:** The process runs as PID 1 inside the namespace, mapped to the host PID shown above.