package utils

import "os"

func CheckDbExistence(dbFilename string) bool {
    if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
        return false
    }

    return true
}

func ReverseBytes(data []byte) {
    for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
        data[i], data[j] = data[j], data[i]
    }
}
