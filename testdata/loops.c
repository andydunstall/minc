fn main() {
	let x = 1;
	loop (x < 5) {
		x = x + 1;
		if (x > 12) {
			break;
		} else if (x > 7) {
			continue;
		}
	}
	return x;
}
