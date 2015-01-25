# go-vcr

The initial release provides the interface and implementation to record your test suite's HTTP interactions and replay them during future test runs for fast, deterministic, accurate tests.

### Usage
```go
import "github.com/ozeias/go-vcr/vcr"
...
server, httpClient := vcr.UseCassette("vine")
client.HTTPClient = httpClient
defer server.Close()

vine, err := client.getVine(vineID)
// ...
```

Once you run this test, go-vcr will record the HTTP request to fixtures/vine.json. When you run it again, go-vcr will replay the response. Works just like [VCR](https://github.com/vcr/vcr).
