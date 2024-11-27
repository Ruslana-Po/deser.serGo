package main

import (
	"fmt"

	"go.ser/hashtable"
)

func main() {
	ht := hashtable.NewHashTable(5)
	ht.AddHash("1", "v1")
	ht.AddHash("2", "v2")
	ht.AddHash("3", "v3")

	fmt.Printf("Binar: \n")
	err := ht.SerializeBinary("hashtable.bin")
	if err != nil {
		fmt.Println("Error:", err)
	}
	newHash := hashtable.NewHashTable(5)
	err = newHash.DeserializeBinary("hashtable.bin")
	if err != nil {
		fmt.Println("Error:", err)
	}
	newHash.Print()

	fmt.Printf("Text: \n")

	err = ht.SerializeText("hashtable.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	newHash2 := hashtable.NewHashTable(5)
	err = newHash2.DeserializeText("hashtable.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}
	newHash2.Print()
}
