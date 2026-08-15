---
title: "Tuning QEMU/KVM Guest Performance on Linux"
slug: b8e2d405
aliases: qemu-kvm-guest-tuning
date: 2026-04-18
tags: [linux, performance, devops]
lang: en
draft: false
---

A default QEMU/KVM virtual machine leaves substantial performance on the table. The hypervisor configuration controls CPU scheduling, memory access latency, network throughput, and disk I/O, and each of these dimensions has practical tuning options that require no kernel patches. This post covers the most impactful changes you can make to a KVM guest on a Linux host.

## Fun Facts

**Fact 1.** KVM became part of the mainline Linux kernel in version 2.6.20, released in February 2007. It was the first hypervisor integrated directly into the Linux kernel source tree.

**Fact 2.** The `virtio` paravirtualized driver family reduces VM exit frequency compared to emulated hardware like e1000 or IDE. Fewer VM exits means fewer expensive context switches between guest and host.

**Fact 3.** Transparent Huge Pages (THP) can actually hurt KVM performance by causing latency spikes during page compaction. Explicit hugepages with `hugetlbfs` provide more predictable latency.

---

## Tips and Tricks

### 1. Pin vCPUs to Physical Cores with taskset and numactl

vCPU threads are ordinary Linux threads. Without pinning, the scheduler can migrate them across NUMA nodes, adding memory access latency.

```bash
# Find the PID of vCPU threads for a running VM
ps -eLf | grep qemu | grep CPU

# Pin vCPU thread (PID 12345) to CPU core 2
taskset -cp 2 12345

# Launch QEMU with CPU affinity from the start
numactl --cpunodebind=0 --membind=0 qemu-system-x86_64 \
  -cpu host \
  -smp 4,sockets=1,cores=4,threads=1 \
  -m 8G \
  -drive file=disk.qcow2,format=qcow2

# Check NUMA topology first
numactl --hardware
```

For production setups, use `libvirt` with `<vcpupin>` and `<numatune>` elements in the domain XML instead of manual `taskset`.

### 2. Configure Hugepages for Lower Memory Latency

Using 2 MB hugepages for guest RAM eliminates TLB pressure and prevents page compaction pauses.

```bash
# Allocate 4096 hugepages of 2 MB each (= 8 GB total)
echo 4096 | sudo tee /proc/sys/vm/nr_hugepages

# Verify allocation
grep HugePages /proc/meminfo

# Make it persistent across reboots
echo "vm.nr_hugepages = 4096" | sudo tee /etc/sysctl.d/99-hugepages.conf

# Mount hugetlbfs if not already mounted
sudo mount -t hugetlbfs hugetlbfs /dev/hugepages
```

Then add `-mem-path /dev/hugepages` to your QEMU command line, or set `<memoryBacking><hugepages/></memoryBacking>` in the libvirt XML.

### 3. Choose virtio-net Over Emulated NICs

The emulated `e1000` or `rtl8139` adapters require QEMU to emulate hardware register I/O for every packet. `virtio-net` uses a shared ring buffer that avoids this overhead.

```bash
# QEMU command line: use virtio-net-pci
qemu-system-x86_64 \
  -netdev tap,id=net0,ifname=tap0,script=no,downscript=no \
  -device virtio-net-pci,netdev=net0,mq=on,vectors=6

# Inside the guest: confirm the driver is loaded
lspci | grep -i ethernet
ethtool -i eth0 | grep driver
# Expected: driver: virtio_net

# Enable multi-queue virtio-net (requires guest support)
ethtool -L eth0 combined 4
```

For maximum throughput, also enable `vhost=on` on the netdev, which moves packet handling into a kernel thread on the host and bypasses QEMU userspace entirely.

### 4. Use io_uring for Disk I/O Passthrough

QEMU's default `aio=threads` disk backend uses a thread pool to emulate async I/O. The `io_uring` backend uses the Linux 5.1+ io_uring interface for lower latency and higher IOPS.

```bash
# Enable io_uring AIO backend
qemu-system-x86_64 \
  -drive file=/var/lib/vms/disk.raw,format=raw,if=none,id=drive0,aio=io_uring,cache=none \
  -device virtio-blk-pci,drive=drive0,num-queues=4

# Benchmark inside the guest with fio
fio --name=randread --ioengine=libaio --iodepth=32 \
    --rw=randread --bs=4k --direct=1 --size=4G \
    --filename=/dev/vda --runtime=30 --time_based

# Compare with aio=threads by swapping the flag and re-running
```

`cache=none` is important here: it disables the host page cache for the disk file, giving io_uring direct access to the storage layer.

### 5. Check and Organize IOMMU Groups for Passthrough

IOMMU groups determine which devices can be passed through together to a guest. Devices in the same group must all be passed through or none of them.

```bash
# List all IOMMU groups and their devices
for g in /sys/kernel/iommu_groups/*; do
  echo "Group $(basename $g):"
  for d in $g/devices/*; do
    echo "  $(basename $d) $(lspci -nns $(basename $d) 2>/dev/null | tail -1)"
  done
done

# Enable IOMMU on Intel systems (add to kernel cmdline in /etc/default/grub)
# GRUB_CMDLINE_LINUX="intel_iommu=on iommu=pt"

# For AMD
# GRUB_CMDLINE_LINUX="amd_iommu=on iommu=pt"

# Verify IOMMU is active
dmesg | grep -i iommu | head -10

# Bind a device to vfio-pci for passthrough
echo "10de 2204" | sudo tee /sys/bus/pci/drivers/vfio-pci/new_id
```

The `iommu=pt` flag enables passthrough mode, which skips IOMMU translation for devices not assigned to guests, reducing overhead for the host.
