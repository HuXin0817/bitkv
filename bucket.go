package bitkv

import (
	"sync"

	"github.com/HuXin0817/bitkv/interval/errors"
)

const headerSize = 2

type bucket struct {
	bits int
	m    map[string]string
	mu   sync.Mutex
}

func newBucket() *bucket {
	return &bucket{
		m: make(map[string]string),
	}
}

func replay(content []byte) (b *bucket, err error) {
	b = newBucket()

	l := len(content)
	// replay log
	if l > 0 {
		var p int
		for p < l {
			kp := p + headerSize
			if kp > l {
				return nil, errors.ErrReplayLog
			}
			vp := kp + int(content[p])
			if vp > l {
				return nil, errors.ErrReplayLog
			}
			end := vp + int(content[p+1])
			if end > l {
				return nil, errors.ErrReplayLog
			}
			k := content[kp:vp]
			v := content[vp:end]
			b.Put(string(k), string(v))
			p = end
		}
	}
	return
}

func (b *bucket) Put(k, v string) (changed bool) {
	if k == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	ov, find := b.m[k]
	if ov == v {
		return false
	}

	if v == "" {
		if find {
			b.bits -= headerSize + len(k) + len(ov)
			delete(b.m, k)
			return true
		}
		return false
	}

	if !find {
		b.bits += headerSize + len(k) + len(v)
		b.m[k] = v
		return true
	}

	b.bits += len(v) - len(ov)
	b.m[k] = v
	return true
}

func (b *bucket) Get(k string) (v string) {
	return b.m[k]
}

func (b *bucket) Export() []byte {
	p := 0
	content := make([]byte, b.bits)

	b.mu.Lock()
	defer b.mu.Unlock()

	for k, v := range b.m {
		keyLength := byte(len(k))
		valueLength := byte(len(v))
		keyPos := p + headerSize
		valuePos := keyPos + int(keyLength)
		endPos := valuePos + int(valueLength)
		content[p] = keyLength
		content[p+1] = valueLength
		copy(content[keyPos:valuePos], k)
		copy(content[valuePos:endPos], v)
		p = endPos
	}

	b.bits = 0
	b.m = make(map[string]string)
	return content
}
