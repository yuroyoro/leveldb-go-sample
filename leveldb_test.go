package main

import (
	"fmt"
	"testing"
)

func TestLevelDB(t *testing.T) {
	dbname := "./testdb"

	db, err := OpenLevelDB(dbname)

	if err != nil {
		fmt.Printf("error : %v\n", err)
		return
	}
	defer db.Close()

	// put keys
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%02d", i)
		val := fmt.Sprintf("val_%02d", i)
		err := db.Put(key, val)
		if err != nil {
			t.Errorf("LevelDB.Put(%s, %s): %v\n", key, val, err)
		}

		fmt.Printf("Put(%s, %s)\n", key, val)
	}

	// get registered keys
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%02d", i)
		val, err := db.Get(key)
		if err != nil {
			t.Errorf("LevelDB.Get(%s): %v\n", key, err)
		}

		expected := fmt.Sprintf("val_%02d", i)
		if val != expected {
			t.Errorf("LevelDB.Get(%s): expected -> %s : actual -> %s\n", key, expected, val)
		}

		fmt.Printf("Get(%s) -> %s\n", key, val)
	}

	// get unknown keys
	for i := 10; i < 12; i++ {
		key := fmt.Sprintf("key_%02d", i)
		val, err := db.Get(key)
		if err != nil {
			t.Errorf("LevelDB.Get(%s): %v\n", key, err)
		}

		if val != "" {
			t.Errorf("LevelDB.Get(%s): should be nil : actual -> %s\n", key, val)
		}

		fmt.Printf("Get(%s) -> %s\n", key, val)
	}

	// delete
	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("key_%02d", i)
		err := db.Delete(key)
		if err != nil {
			t.Errorf("LevelDB.Delete(%s): %v\n", key, err)
		}

		fmt.Printf("Delete(%s)\n", key)
	}
	// get deleted keys
	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("key_%02d", i)
		val, err := db.Get(key)
		if err != nil {
			t.Errorf("LevelDB.Get(%s): %v\n", key, err)
		}
		if val != "" {
			t.Errorf("LevelDB.Get(%s): should be deleted: actual -> %s\n", key, val)
		}

		fmt.Printf("Get(%s) -> %s\n", key, val)
	}

}
