package main

func BuildRunner() Runner {
	c := Runner{}
	c.Inject(Store{}, &InMemoryBackend{})
	return c
}

func ExampleVersion() {
	runner := BuildRunner()
	runner.Delegate("version", "", "")
	// Output: 0.0.1
}
