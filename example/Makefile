compile: tester.cpp student.cpp correct.cpp
	mkdir -p out
	clang++ -std=c++20 -o out/tester tester.cpp student.cpp
	clang++ -std=c++20 -o out/comparison tester.cpp correct.cpp

test: compile
	go-vpl-tester --timeout=100ms --entrypoint=./out/tester --max-grade=2