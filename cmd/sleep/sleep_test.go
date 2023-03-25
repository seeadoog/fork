package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
	"unsafe"
)

func TestName1(t *testing.T) {

	a := 5
	{
		a := 3
		fmt.Println(a)
	}
	fmt.Println(a)
}

func TestName2(t *testing.T) {
	in := [3]string{"a", "b", "c"}
	var out []*string
	for _, v := range in {
		//copy := v
		v := v
		out = append(out, &v)
	}
	fmt.Println("Values:", *out[0], *out[1], *out[2])

}

func TestName(t *testing.T) {
	var a int = 5
	ptr := &a
	fmt.Println(*ptr)
	a = 15
	fmt.Println(*ptr)
}

//2020-05-16 19:20:34|user.login|name=Charles&location=Beijing&device=iPhone

func parseAndFormatString(in string) (out string) {
	res := make(map[string]interface{})
	for _, ss := range strings.Split(in, "|") {
		for _, s := range strings.Split(ss, "&") {
			key, val, ok := strings.Cut(s, "=")
			if ok {
				res[key] = val
			}
		}
	}
	bs, _ := json.Marshal(res)
	bf := bytes.NewBuffer(nil)
	bf.Grow(len(bs))
	json.Indent(bf, bs, "", "\t")
	bs = bf.Bytes()
	return *(*string)(unsafe.Pointer(&bs))
}

func TestJSON(t *testing.T) {
	fmt.Println(parseAndFormatString("2020-05-16 19:20:34|user.login|name=Charles&location=Beijing&device=iPhone"))
}

// -10 ,-8,0,1,5,7
// -10 ,-8,-2,3,5,7
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
func findMinAbs(arr []int) int {
	lo := 0
	hi := len(arr) - 1
	mid := 0
	for lo <= hi {
		mid = (lo + hi) / 2
		v := arr[mid]
		if v > 0 {
			if mid > 0 && arr[mid-1] < 0 {
				if v < -arr[mid-1] {
					return v
				} else {
					return arr[mid-1]
				}
			}
			hi = mid - 1
		} else if v == 0 {
			return v
		} else if v < 0 {
			if mid < len(arr)-1 && arr[mid+1] > 0 {
				if -v < arr[mid+1] {
					return v
				} else {
					return arr[mid+1]
				}

			}
			lo = mid + 1
		}
	}
	return arr[mid]
}

func TestFindABs(t *testing.T) {

	fmt.Println(findMinAbs([]int{-10, -5, 1, 2, 5}))
	fmt.Println(findMinAbs([]int{-10, -2, 3, 4, 5}))
	fmt.Println(findMinAbs([]int{-10, -2, 0, 3, 4, 5}))
	fmt.Println(findMinAbs([]int{1, 2, 3, 4, 5}))
	fmt.Println(findMinAbs([]int{-1, 2, 3, 4, 5}))
	fmt.Println(findMinAbs([]int{-3, 2, 3, 4, 5}))
	fmt.Println(findMinAbs([]int{-10, -5, -3, -2}))
}

func sortArray(nums []int) []int {
	return quicksort(nums, 0, len(nums)-1)
}

func quicksort(nums []int, low int, high int) []int {
	if low >= high {
		return nums
	}
	pivot := nums[low]
	start := low
	end := high
	for low < high {
		for low < high && nums[high] >= pivot {
			high--
		}
		if low < high {
			nums[low] = nums[high]
			low++
		}
		for low < high && nums[low] < pivot {
			low++
		}
		if low < high {
			nums[high] = nums[low]
			high--
		}
	}
	nums[low] = pivot
	quicksort(nums, start, low-1)
	quicksort(nums, low+1, end)
	return nums
}

func findMid(arr []int) int {
	sortArray(arr)
	return arr[len(arr)/2]
}

func TestQW(t *testing.T) {
	arr := []int{8, 6, 5, 7, 3, 4, 2, 10, 1}
	//quickSort1(arr, 0, len(arr))
	sortArray(arr)
	fmt.Println(arr)
}

func nStep(n int) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	}
	arr := make([]int, n+1)
	arr[0] = 0
	arr[1] = 1
	arr[2] = 2
	for i := 3; i < len(arr); i++ {
		arr[i] = arr[i-1] + arr[i-2]
	}
	return arr[n]
}

func TestNN(t *testing.T) {

	fmt.Println(nStep(5))
}
func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func qsort(arr []int, lo int, hi int) {
	if hi-lo < 1 {
		return
	}
	i := lo + 1
	j := hi
	p := arr[lo]
	for {
		for i <= j && arr[i] <= p {
			i++
		}
		for i <= j && arr[j] >= p {
			j--
		}
		if i >= j {
			break
		}
		swap(arr, i, j)
	}
	if p > arr[i-1] {
		swap(arr, lo, i-1)
	}
	qsort(arr, lo, i-2)
	qsort(arr, i, hi)
}

func TestQs(t *testing.T) {
	arr := []int{1, 11, 4, 5, 3, 7, 9, 2, 8, 19}
	//arr := []int{3, 2, 4, 7, 9, 5, 8}
	qsort(arr, 0, len(arr)-1)
	fmt.Println(arr)

}

func TestWDSD(t *testing.T) {
	a := make([]int, 100300)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < len(a); i++ {
		j := rand.Int() % len(a)
		swap(a, i, j)
	}
	fmt.Println(sort.SliceIsSorted(a, func(i, j int) bool {
		return a[i] < a[j]
	}))

	qsort(a, 0, len(a)-1)

	fmt.Println(sort.SliceIsSorted(a, func(i, j int) bool {
		return a[i] < a[j]
	}))

}
