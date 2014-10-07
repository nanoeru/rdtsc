//	RDTSC(Read Time Stamp Counter)
package rdtsc

import (
	"fmt"
	"time"
)

type Cycle struct {
	Hi uint32
	Lo uint32
}

//	in asm
func rdtsc(hi *uint32, lo *uint32)

//	Cycle Counter 取得
func GetCounterHiLo() (cycHi, cycLo uint32) {
	rdtsc(&cycHi, &cycLo)
	return
}

//	Cycle Counter 取得
func GetCounter() (cycle Cycle) {
	rdtsc(&cycle.Hi, &cycle.Lo)
	return
}

//	Cycle Counter 計算 (現在を基準)
func CalNowCounterHiLo(cycHi, cycLo uint32) float64 {
	var ncycHi, ncycLo uint32
	rdtsc(&ncycHi, &ncycLo)
	return CalCounterHiLo(cycHi, cycLo, ncycHi, ncycLo)
}

//	Cycle Counter 計算 (現在を基準)
func CalNowCounter(cycle Cycle) float64 {
	ncycle := GetCounter()
	return CalCounter(cycle, ncycle)
}

//	Cycle Counter 計算 (before・after)
func CalCounter(cycle Cycle, ncycle Cycle) (result float64) {
	return CalCounterHiLo(cycle.Hi, cycle.Lo, ncycle.Hi, ncycle.Lo)
}

//	Cycle Counter 計算 (before・after)
func CalCounterHiLo(cycHi, cycLo, ncycHi, ncycLo uint32) (result float64) {
	var hi, lo, borrow uint32
	// Do double precision subtraction
	lo = ncycLo - cycLo
	if lo > ncycLo {
		borrow = 1
	} else {
		borrow = 0
	}
	hi = ncycHi - cycHi - borrow

	result = float64(hi)*(1<<30)*4 + float64(lo)
	if result < 0.0 {
		//fmt.Fprintf(os.Stderr, "Error: Cycle counter returning negative value: %.0f\n", result)
		panic(fmt.Sprintf("Error: Cycle counter returning negative value: %.0f\n", result))
	}
	return result
}

//	クロック取得(フル設定)
func MhzFull(sleepTime int, verbose bool) (rate float64) {
	cycHi, cycLo := GetCounterHiLo()
	time.Sleep(time.Second * time.Duration(sleepTime))
	cycs := CalNowCounterHiLo(cycHi, cycLo)
	rate = cycs / (1e6 * float64(sleepTime))
	if verbose {
		fmt.Printf("Sleeping %d seconds required %.0f clock cycles\n", sleepTime, cycs)
		fmt.Printf("Processor Clock Rate ~= %.1f MHz\n", rate)
	}
	return
}

// クロック取得
func Mhz(verbose bool) float64 {
	return MhzFull(1, verbose)
}
