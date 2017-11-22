package registry

type mockRegistry struct {
	Services map[string][]*Service
}

var (
	mockData = map[string][]*Service{
		"foo": []*Service{
			{
				Name:    "foo",
				Version: "1.0.0",
				Nodes: []*Node{
					{
						Id:      "foo-1.0.0-123",
						Address: "localhost",
						Port:    9999,
					},
					{
						Id:      "foo-1.0.0-321",
						Address: "localhost",
						Port:    9999,
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.1",
				Nodes: []*Node{
					{
						Id:      "foo-1.0.1-321",
						Address: "localhost",
						Port:    6666,
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.3",
				Nodes: []*Node{
					{
						Id:      "foo-1.0.3-345",
						Address: "localhost",
						Port:    8888,
					},
				},
			},
		},
	}
)

func (m *mockRegistry) init() {
	// add some mock data
	m.Services = mockData
}

func (m *mockRegistry) GetService(service string) ([]*Service, error) {
	s, ok := m.Services[service]
	if !ok || len(s) == 0 {
		return nil, ErrNotFound
	}
	return s, nil

}

func (m *mockRegistry) ListServices() ([]*Service, error) {
	var services []*Service
	for _, service := range m.Services {
		services = append(services, service...)
	}
	return services, nil
}

func (m *mockRegistry) Register(s *Service, opts ...RegisterOption) error {
	services := addServices(m.Services[s.Name], []*Service{s})
	m.Services[s.Name] = services
	return nil
}

func (m *mockRegistry) Deregister(s *Service) error {
	services := delServices(m.Services[s.Name], []*Service{s})
	m.Services[s.Name] = services
	return nil
}

func (m *mockRegistry) Watch() (Watcher, error) {
	return &mockWatcher{exit: make(chan bool)}, nil
}

func (m *mockRegistry) String() string {
	return "mock"
}

func NewMockRegistry() Registry {
	m := &mockRegistry{Services: make(map[string][]*Service)}
	m.init()
	return m
}
