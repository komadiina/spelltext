// package audio

// import (
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/gopxl/beep"
// 	"github.com/gopxl/beep/mp3"
// 	"github.com/gopxl/beep/speaker"
// )

// var initOnce sync.Once

// func initSpeaker(sr beep.SampleRate) {
// 	initOnce.Do(func() {
// 		speaker.Init(sr, sr.N(time.Second/10))
// 	})
// }

// var sounds map[string]beep.Buffer

// func PlaySound(path string) error {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}

// 	streamer, format, err := mp3.Decode(f)
// 	if err != nil {
// 		f.Close()
// 		return err
// 	}
// 	initSpeaker(format.SampleRate)

// 	seq := beep.Seq(streamer, beep.Callback(func() {
// 		streamer.Close()
// 		f.Close()
// 	}))

// 	speaker.Play(seq)
// 	return nil
// }

package audio

import (
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/utils/singleton/logging"
)

type Manager struct {
	AudioEnabled bool
	targetRate   beep.SampleRate
	buffers      map[string]*beep.Buffer
	mutex        sync.RWMutex
	initOnce     sync.Once
}

func NewManager(targetRate beep.SampleRate) *Manager {
	return &Manager{
		targetRate: targetRate,
		buffers:    make(map[string]*beep.Buffer),
	}
}

func (m *Manager) Preload(soundFiles []string) error {
	m.initOnce.Do(func() {
		speaker.Init(m.targetRate, m.targetRate.N(time.Second/10))
	})

	for _, key := range soundFiles {
		f, err := os.Open(key)
		if err != nil {
			return err
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			f.Close()
			return err
		}

		var source beep.Streamer = streamer
		buffer := beep.NewBuffer(format)
		buffer.Append(source)

		streamer.Close()
		f.Close()

		m.mutex.Lock()
		m.buffers[key] = buffer
		m.mutex.Unlock()
	}

	return nil
}

func (m *Manager) Play(key string, logger *logging.Logger) {
	if !m.AudioEnabled {
		return
	}

	m.mutex.RLock()
	buf, ok := m.buffers[key]
	m.mutex.RUnlock()

	if !ok {
		logger.Warnf("sound %s not found", key)
		return
	}

	s := buf.Streamer(0, buf.Len())

	speaker.Play(s)
}

func (m *Manager) PlayBackground(logger *logging.Logger) {
	if !m.AudioEnabled {
		return
	}

	m.mutex.RLock()
	// use beep.Loop(-1, streamer) for constants.BACKGROUND_OST
	buf, ok := m.buffers[constants.BACKGROUND_OST]
	m.mutex.RUnlock()

	if !ok {
		logger.Warnf("sound %s not found", constants.BACKGROUND_OST)
		return
	}

	str := buf.Streamer(0, buf.Len())
	loop := beep.Loop(-1, str)
	speaker.Play(loop)
}
