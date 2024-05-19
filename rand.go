package mygo

import (
	"math/rand"
	"time"
)

func (*GoRand) String(length int) string {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	time.Sleep(10 * time.Nanosecond)

	letter := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// Ints 获取指定范围内的随机序列
func (*GoRand) Ints(from, to, size int) []int {
	if to-from < size {
		size = to - from
	}

	var slice []int
	for i := from; i < to; i++ {
		slice = append(slice, i)
	}

	var ret []int
	for i := 0; i < size; i++ {
		idx := rand.Intn(len(slice))
		ret = append(ret, slice[idx])
		slice = append(slice[:idx], slice[idx+1:]...) // 剔除已经选中的
	}
	return ret
}

// Int 随机获取指定范围内的某一个数
func (*GoRand) Int(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(10 * time.Nanosecond)
	return min + rand.Intn(max-min)
}
