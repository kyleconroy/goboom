package main

func ExampleVersion() {
	runner := Runner{storage: Store{}, backend: &InMemoryBackend{}}
	runner.Delegate("version", "", "")
	// Output: 0.0.1
}
