# Minc

Minc is a simple programming language, which aims to be a simplified version of
C, combining ideas from C, Go and Rust.

## Variables

Minc variables are similar to Rust:
```
let a: u32 = 100;
let mut b: i32 = 300;
```

* Integer types:
  * `i8`/`u8`
  * `i16`/`u16`
  * `i32`/`u32`
  * ...
* Variables are immutable by default

## Functions

Minc function definitions are similar tuo Rust:
```
fn inc(a: u32) -> u32 {
	return a + 1;
}
```
