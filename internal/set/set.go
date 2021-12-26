package set

type String struct {
	value map[string]bool
}

func NewString(values ...string) *String {
	s := &String{
		value: make(map[string]bool),
	}
	s.Add(values...)
	return s
}

func (s *String) Len() int {
	return len(s.value)
}

func (s *String) Add(items ...string) {
	for _, item := range items {
		s.value[item] = true
	}
}

func (s *String) Remove(items ...string) {
	for _, item := range items {
		delete(s.value, item)
	}
}

func (s *String) Has(item string) bool {
	_, exists := s.value[item]
	return exists
}

func (s *String) Intersect(other *String) *String {
	if s.Len() > other.Len() {
		return other.Intersect(s)
	}
	ns := NewString()
	for k := range s.value {
		if other.Has(k) {
			ns.Add(k)
		}
	}
	return ns
}

func (s *String) Difference(other *String) *String {
	ns := NewString(s.List()...)
	for _, k := range other.List() {
		ns.Remove(k)
	}
	return ns
}

func (s *String) Union(other *String) *String {
	ns := NewString()
	ns.Add(s.List()...)
	ns.Add(other.List()...)
	return ns
}

func (s *String) Equal(other *String) bool {
	for k := range s.value {
		if !other.Has(k) {
			return false
		}
	}
	return s.Len() == other.Len()
}

func (s *String) List() []string {
	val := make([]string, 0, len(s.value))
	for k := range s.value {
		val = append(val, k)
	}
	return val
}

func (s *String) Map() map[string]bool {
	m := make(map[string]bool)
	for k := range s.value {
		m[k] = true
	}
	return m
}
