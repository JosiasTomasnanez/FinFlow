package lab

import "sync"

// Hold retiene un bloque en el proceso para que el working set de cAdvisor
// no baje al terminar el HTTP handler (el GC se lo comería).
type Hold struct {
	mu   sync.Mutex
	blob []byte
}

func touchPages(buf []byte) {
	for i := 0; i < len(buf); i += 4096 {
		buf[i] = 1
	}
}

func (h *Hold) Retain(mebibytes int) int {
	mebibytes = ClampInt(mebibytes, DefaultMemoryMiB, MaxMemoryMiB)
	size := mebibytes << 20

	h.mu.Lock()
	if cap(h.blob) >= size {
		h.blob = h.blob[:size]
		touchPages(h.blob)
		h.mu.Unlock()
		return mebibytes
	}
	h.blob = nil
	h.mu.Unlock()

	buf := make([]byte, size)
	touchPages(buf)

	h.mu.Lock()
	h.blob = buf
	h.mu.Unlock()
	return mebibytes
}

func (h *Hold) Release() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.blob = nil
}

func (h *Hold) MiB() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.blob) >> 20
}
