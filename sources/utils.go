package main

import "math/rand"

func generateRandomArray(arraylength, stringLength int) []string {
	res := make([]string, arraylength)
	for i := range res {
		res[i] = generateRandomString(stringLength)
	}

	return res
}

func generateRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
