# jitter

Reusable jitter strategies to use for backoff. Minimal API, pluggable strategies, production-ready.

This package provides a small set of jitter/backoff strategies based on the AWS ["Exponential
Backoff and Jitter"](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

## Install

```sh
go get github.com/khinshankhan/jitter-go
```

## Overview

The core API is intentionally tiny:

```go
type Strategy interface {
    Next(attempt int) int64
}

type Config struct {
    Base   int64
    Cap    int64
    Random func(max int64) int64
}
```

You choose the time units (ms, seconds, etc.) and convert them however you like.

## Example

Here's the simplest way to use it:

```go
package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/khinshankhan/jitter-go"
)

func main() {
	ctx := context.Background()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	backoff := jitter.New(jitter.Config{
		Base:   100,      // 100 ms units
		Cap:    10_000,   // cap at 10s (in ms units)
		Random: r.Int63n, // U[0, n)
	})

	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delayMs := backoff.Next(attempt)
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
		if err := doSomething(ctx); err == nil {
			break
		}
	}

	// all retries failed
}

func doSomething(ctx context.Context) error {
	// your operation here
	return nil
}
```

### A safer, cancellable pattern

I highly recommend using it with `context` in a `select` statement.

Your actual function likely takes in a context, and the business logic may include deadlines or
timeouts (eg Lambdas, Cloud Run tasks, cronjobs, etc). This is how it may look:

```go
r := rand.New(rand.NewSource(time.Now().UnixNano()))
backoff := jitter.New(jitter.Config{
	Base:   100,
	Cap:    10_000,
	Random: r.Int63n,
})

const maxAttempts = 5
for attempt := 0; attempt < maxAttempts; attempt++ {
	if attempt > 0 {
		delayMs := backoff.Next(attempt)
		delay := time.Duration(delayMs) * time.Millisecond

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// continue to the retry
		}
	}

	if err := doSomething(ctx); err == nil {
		return nil
	}
}

// all retries failed
return fmt.Errorf("retries exhausted")
```

This makes cancellation immediate and avoids the temporary leaks associated with `time.After` (more
below).

## Notes

This package started as an exercise and a minimal abstraction of the standard jitter strategies
described in the AWS blogpost [Exponential Backoff And
Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/).

These algorithms are widely used -- even the Google SRE book references the AWS article under
[Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/).

A surprising thing I noticed while exploring existing implementations (AWS SDKs, GCP libraries,
[cenkalti/backoff](https://github.com/cenkalti/backoff), etc) is that they're usually "batteries
included": lots of config knobs, but you rarely get to learn or swap the actual jitter
algorithms. This package intentionally separates the strategies so you can explore, compare, or
extend them.

I learned a fair bit implementing this recently and decided to extract it into a package as a mini
blogpost of sorts of my learnings.

### Pseudo Random

This topic is a rabbit hole.

#### Rabbit Hole

Out of the box, `math/rand` with a shared Rand is not concurrency safe. Go `math/rand` does
implement top level random functions which _are_ safe for concurrency but it seems to use a simple
mutex.

However, when using multiple threads, mutexes are the bane of performance no matter what language
you use. It would be _correct_ but locking is not _ideal_. Using a platform specific primitive like
`channels` in go or `CompletableFuture` in Java will always be better optimized by the compiler and
runtime. However, this is a much more complicated...

In theory you could have a per-thread PRNG seed which makes things better and initialize it from
real randomness when it's first called in a thread. Of course if multiple threads read the seed
before it is updated they all get the same next random number. So you either explicitly thread the
prng seed through all your calls (inconvenient) or you put the seed in thread-local storage, and
then your threads are all doing independent things. But if you're spawning and destroying lots of
threads and each of them is calling rand, then you may be pulling a lot of entropy and you may run
out...

The simplest way I've found is to simply read from the random source, which could very well be
`time.Now().UnixNano()`, add a number unique to each routine, and then use that as your seed per
routine and keep the routines bounded (don't kill and create them constantly).

#### Solution

In the above rabbit hole, that is _my_ understanding and a solution good enough for the initial
project I've made this package for. However, it should be recognized that others may have other
solutions or other thoughts.

This is why I intentionally set it up so we intake `func(max int64) int64`, which lets you plug in
whatever fits your workload or even an evolution in understanding in tech. 0 opinions.

### Waiting

`time.Sleep` doesn't have a guarantee system will have enough enough resources to wake
up. `time.Sleep(d)` basically says "sleep at least d", not "wake up exactly at d". It depends on:

- Go's scheduler
- the OS scheduler
- system load / CPU availability

But this is true no matter what you use (`Sleep`, `After`, `NewTimer`). The underlying limitation is
the OS + runtime, not the specific API.

For retry/backoff, "approximately this long" is usually totally fine, you'll just end up waiting a
little longer in the worst case. If you need hard real-time guarantees, Go timers won't get you that
anyway.

`time.After(d)` allocates a timer that you **cannot cancel early**, the runtime only releases it
when the timer fires, creating a _temporary leak_. In long-running loops, this can create temporary
buildup.

> [!TIP]
>
> This isn't relevant anymore from Go 1.23
>
> Read more about it https://go.dev/wiki/Go123Timer but basically they changed the implementation
> :slightly_smiling_face:

`time.NewTimer(d)` lets you:

- cancel early if another condition fires first
- avoid the `time.After` leak (you can freely cancel before time via `Stop`)
- integrate cleanly with `select` + `context`

This is why the second example above uses `NewTimer`.

## Contributing

Bug reports and PRs are welcome. Please open an issue first for discussion.

Feel free to open an issue if you spot something iffy or have a hot tip :shrug:

## License

Apache-2.0. See [LICENSE](./LICENSE). If applicable, see [NOTICE](./NOTICE).
