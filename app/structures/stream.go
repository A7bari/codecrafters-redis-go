package structures

type Entry struct {
	Key   string
	Pairs map[string]string
}

type Stream struct {
	Entries []Entry
}

func NewStream() *Stream {
	return &Stream{Entries: []Entry{}}
}

func (s *Stream) Add(key string, pairs map[string]string) {
	s.Entries = append(s.Entries, Entry{Key: key, Pairs: pairs})
}

func (s *Stream) Get(key string) (map[string]string, bool) {
	for _, e := range s.Entries {
		if e.Key == key {
			return e.Pairs, true
		}
	}
	return nil, false
}
