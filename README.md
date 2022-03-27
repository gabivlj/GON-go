# GON in Go

File format created by Tyler Glaiel ([source](https://github.com/TylerGlaiel/GON)).

## What it does

- At the moment it only tries to deserialize into a hashmap a single GON file, it doesn't have all the cool stuff that the original repository has.
- It's kind of fast because I try to not make many allocations and keep it simple, files that aren't well formed will fail catastrophically.
- Not a serious implementation, it's only for fun.
- Supports arrays, objects, strings, floats, integers and booleans.
