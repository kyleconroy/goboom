package main

func ExampleVersion() {
	runner := Runner{storage: Store{}, backend: &InMemoryBackend{}}
	runner.Delegate("version", "", "")
	// Output: 0.0.1
}

func ExampleEcho() {
	store := Store{}
	store["foo"] = map[string]string{"bar": "bat"}

	runner := Runner{storage: store, backend: &InMemoryBackend{}}
	runner.Delegate("echo", "foo", "bar")
	// Output: bat
}
