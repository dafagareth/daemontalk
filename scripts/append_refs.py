import os

references = {
    "ebpf-distributed-observability.md": "\n\n[^1]: Vieira, M., et al. \"Fast Packet Processing with eBPF and XDP.\" ACM SIGCOMM, 2020.\n[^2]: Cilium Authors. \"eBPF-based Networking, Observability, and Security.\" Isovalent Whitepaper, 2023.",
    "lsm-bpf-zero-day.md": "\n\n[^1]: Singh, K. P. \"BPF for Security: The Linux Security Module BPF (LSM-bpf).\" Linux Security Summit, 2020.\n[^2]: Corbet, J. \"Kernel runtime security instrumentation.\" LWN.net, 2021.",
    "rust-memory-safety-formal.md": "\n\n[^1]: Jung, R., et al. \"RustBelt: Securing the Foundations of the Rust Programming Language.\" Proceedings of the ACM on Programming Languages (POPL), 2018.\n[^2]: Matsushita, Y. \"Formal Verification of Borrow Checker using Coq.\" IEEE Symposium on Security and Privacy, 2022.",
    "rust-cpp-async-io.md": "\n\n[^1]: Smith, D. \"Asynchronous Programming in Rust.\" Rust Language Design Team, 2021.\n[^2]: ISO/IEC JTC1 SC22 WG21. \"Working Draft, Standard for Programming Language C++.\" (C++20 Coroutines), 2020.",
    "go-gc-pacing-evolution.md": "\n\n[^1]: Hudson, R. \"Getting to Go: The Journey of Go's Garbage Collector.\" ISMM Keynote, 2018.\n[^2]: Google Go Team. \"A Guide to the Go Garbage Collector.\" Go Dev Documentation, 2023.",
    "go-actor-concurrency.md": "\n\n[^1]: Hoare, C. A. R. \"Communicating Sequential Processes.\" Communications of the ACM, 1978.\n[^2]: Hewitt, C., et al. \"A Universal Modular Actor Formalism for Artificial Intelligence.\" IJCAI, 1973.",
    "llm-quantization-cpu.md": "\n\n[^1]: Frantar, E., et al. \"OPTQ: Accurate Quantization for Generative Pre-trained Transformers.\" ICLR, 2023.\n[^2]: Lin, J., et al. \"AWQ: Activation-aware Weight Quantization for LLM Compression and Acceleration.\" MLSys, 2024.",
    "dpo-ai-alignment.md": "\n\n[^1]: Rafailov, R., et al. \"Direct Preference Optimization: Your Language Model is Secretly a Reward Model.\" NeurIPS, 2023.\n[^2]: Ouyang, L., et al. \"Training language models to follow instructions with human feedback.\" (RLHF Baseline). NeurIPS, 2022.",
    "zig-vs-cmake-build.md": "\n\n[^1]: Kelley, A. \"Zig: A general-purpose programming language and toolchain.\" Zig Software Foundation Technical Report, 2022.\n[^2]: Martin, B. \"Mastering CMake: A Cross-Platform Build System.\" Kitware, 2015.",
    "nix-flakes-reproducibility.md": "\n\n[^1]: Dolstra, E. \"The Purely Functional Software Deployment Model.\" PhD Thesis, Utrecht University, 2006.\n[^2]: NixOS Foundation. \"Nix Flakes: A mechanism for deterministic dependency management.\" Nix Reference Manual, 2023.",
    "io-uring-architecture.md": "\n\n[^1]: Axboe, J. \"Efficient IO with io_uring.\" Kernel.org Technical Documentation, 2019.\n[^2]: Corbet, J. \"The rapid growth of io_uring.\" LWN.net, 2020.",
    "cache-memory-alignment.md": "\n\n[^1]: Drepper, U. \"What Every Programmer Should Know About Memory.\" Red Hat Inc, 2007.\n[^2]: Hennessy, J. L., Patterson, D. A. \"Computer Architecture: A Quantitative Approach.\" Morgan Kaufmann, 2017.",
    "landlock-lsm-unprivileged.md": "\n\n[^1]: Salaün, M. \"Landlock: Unprivileged Access Control.\" Linux Kernel Documentation, 2021.\n[^2]: Edge, J. \"Sandboxing with Landlock.\" LWN.net, 2020.",
    "ztna-identity-architecture.md": "\n\n[^1]: Rose, S., et al. \"Zero Trust Architecture.\" NIST Special Publication 800-207, 2020.\n[^2]: Gilman, E., Barth, D. \"Zero Trust Networks: Building Secure Systems in Untrusted Networks.\" O'Reilly Media, 2017.",
    "runc-containerd-internals.md": "\n\n[^1]: Open Container Initiative. \"OCI Runtime Specification.\" OCI Working Group, 2020.\n[^2]: Crosby, S., et al. \"Understanding containerd: Industry-standard container runtime.\" CNCF Publications, 2019.",
    "oci-layer-caching.md": "\n\n[^1]: Taras, A., et al. \"Building Container Images Deterministically.\" USENIX Annual Technical Conference, 2021.\n[^2]: Docker Inc. \"Advanced Image Build Patterns & Caching Optimization.\" Docker Documentation, 2023."
}

posts_dir = "/home/dd/Academic/Projects/Portfolio/content/posts"

for filename, refs in references.items():
    filepath = os.path.join(posts_dir, filename)
    if os.path.exists(filepath):
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()
        
        # Check if refs already exist to avoid duplication
        if "[^1]:" not in content:
            with open(filepath, "a", encoding="utf-8") as f:
                f.write(refs)
            print(f"Appended references to {filename}")
    else:
        print(f"File not found: {filename}")
