package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano())) // один источник
	mu    sync.Mutex                                        // защита от гонок при одновременном доступе
	chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
)

// NewRandomString генерирует случайную строку заданной длины
func NewRandomString(size int) string {
	b := make([]rune, size)

	mu.Lock() // блокируем доступ, если функция вызывается из разных горутин
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	mu.Unlock()

	return string(b)
}
