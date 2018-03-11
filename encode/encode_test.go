package encode

import "testing"

func TestEncode(t *testing.T){
	pass := "angryMonkey"
	hash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	hashTest := Encode(pass)
	if (hashTest != hash){
		t.Errorf("Error hashing %v", pass)
		t.Errorf("Result: %v", hashTest)
		t.Errorf("Correct: %v", hash)
	}
}