fn two() {
	return 2;
}

fn addFive(int n) {
	return n + 5;
}

fn addTen(int n) {
	return addFive(addFive(n));
}

fn main() {
	return two() + addTen(two() + 1) + addTen(5);
}
