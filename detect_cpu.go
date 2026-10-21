package useragent

func detectCPU(ua *UserAgent, ts *tokenSet) {
	if ua.CPU.Architecture != "" {
		return
	}
	switch {
	case ts.has("arm64"), ts.has("aarch64"):
		ua.CPU.Architecture = "arm64"
	case ts.has("x86_64"), ts.has("amd64"),
		ts.has("Win64"), ts.has("x64"):
		ua.CPU.Architecture = "amd64"
	case ts.has("WOW64"):
		ua.CPU.Architecture = "386"
	case ts.has("i686"), ts.has("i386"), ts.has("i86pc"):
		ua.CPU.Architecture = "386"
	case ts.has("armv7l"), ts.has("armv8l"), ts.has("armv7"):
		ua.CPU.Architecture = "arm"
	case ts.has("ppc64"):
		ua.CPU.Architecture = "ppc64"
	case ts.has("riscv64"):
		ua.CPU.Architecture = "riscv64"
	case ts.has("mips64"), ts.has("mips"):
		ua.CPU.Architecture = "mips"
	}
}
