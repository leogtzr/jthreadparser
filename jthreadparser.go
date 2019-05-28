package jthreadparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	threadInformationBegins     = "\""
	threadNameRgx               = `^\"(.*)\".*prio=([0-9]+) tid=(\w*) nid=(\w*)\s\w*`
	stateRgx                    = `\s+java.lang.Thread.State: (.*)`
	lockedRgx                   = `\s*\- locked\s*<(.*)>\s*\(a\s(.*)\)`
	parkingOrWaitingRgx         = `\s*\- (?:waiting on|parking to wait for)\s*<(.*)>\s*\(a\s(.*)\)`
	stackTraceRgx               = `^\s+(at|\-\s).*\)$`
	threadNameRgxGroupIndex     = 1
	threadPriorityRgxGroupIndex = 2
	threadIDRgxGroupIndex       = 3
	threadNativeIDRgxGroupIndex = 4
)

// ThreadInfo ...
type ThreadInfo struct {
	Name, ID, NativeID, Priority, State, StackTrace string
	Daemon                                          bool
}

// Locked ...
type Locked struct {
	LockID, LockecObjectName string
}

func (th ThreadInfo) String() string {
	if th.Daemon {
		return fmt.Sprintf("Thread Id: '%s' (daemon), Name: '%s', State: '%s'", th.ID, th.Name, th.State)
	}
	return fmt.Sprintf("Thread Id: '%s', Name: '%s', State: '%s'", th.ID, th.Name, th.State)
}

func extractThreadInfoFromLine(line string) ThreadInfo {
	ti := ThreadInfo{}
	if rgxp, _ := regexp.Compile(threadNameRgx); rgxp.MatchString(line) {
		if strings.Contains(line, " daemon ") {
			ti.Daemon = true
		}
		for _, v := range rgxp.FindAllStringSubmatch(line, -1) {
			ti.Name, ti.Priority, ti.ID, ti.NativeID =
				v[threadNameRgxGroupIndex], v[threadPriorityRgxGroupIndex], v[threadIDRgxGroupIndex], v[threadNativeIDRgxGroupIndex]
		}
	}
	return ti
}

func extractThreadState(line string) string {
	if rgxp, _ := regexp.Compile(stateRgx); rgxp.MatchString(line) {
		st := strings.Split(line, ":")
		if st != nil {
			state := strings.TrimSpace(st[1])
			return strings.Split(state, " ")[0]
		}
	}
	return ""
}

func extractThreadStackTrace(scanner *bufio.Scanner) string {
	var buffer bytes.Buffer
	for len(strings.TrimSpace(scanner.Text())) != 0 {
		buffer.WriteString(strings.TrimSpace(scanner.Text()))
		buffer.WriteString("\n")
		scanner.Scan()
	}
	return buffer.String()
}

// ParseFromFile ...
func ParseFromFile(fileName string) ([]ThreadInfo, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	threads := make([]ThreadInfo, 0)
	parse(file, &threads)

	return threads, nil
}

func parse(r io.Reader, threads *[]ThreadInfo) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), threadInformationBegins) {
			ti := extractThreadInfoFromLine(scanner.Text())
			scanner.Scan()
			ti.State = extractThreadState(scanner.Text())
			scanner.Scan()
			ti.StackTrace = extractThreadStackTrace(scanner)
			*threads = append(*threads, ti)
		}

	}
}

// ParseFrom ...
func ParseFrom(r io.Reader) ([]ThreadInfo, error) {
	threads := make([]ThreadInfo, 0)
	parse(r, &threads)
	return threads, nil
}

// Holds ...
func Holds(threads *[]ThreadInfo) map[ThreadInfo][]Locked {

	holds := make(map[ThreadInfo][]Locked)

	for _, th := range *threads {
		if len(th.StackTrace) == 0 || !strings.Contains(th.StackTrace, "locked") {
			continue
		}

		for _, stackLine := range strings.Split(th.StackTrace, "\n") {
			if rgxp, _ := regexp.Compile(lockedRgx); rgxp.MatchString(stackLine) {
				for _, group := range rgxp.FindAllStringSubmatch(stackLine, -1) {
					if _, exists := holds[th]; !exists {
						holds[th] = make([]Locked, 0)
					}
					h := holds[th]
					h = append(h, Locked{group[1], group[2]})
					holds[th] = h
				}
			}
		}
	}

	return holds
}

// HoldsForThread ...
func HoldsForThread(th *ThreadInfo) []Locked {
	holds := make([]Locked, 0)

	if len(th.StackTrace) == 0 || !strings.Contains(th.StackTrace, "locked") {
		return []Locked{}
	}

	for _, stackLine := range strings.Split(th.StackTrace, "\n") {
		if rgxp, _ := regexp.Compile(lockedRgx); rgxp.MatchString(stackLine) {
			for _, group := range rgxp.FindAllStringSubmatch(stackLine, -1) {
				lock := Locked{group[1], group[2]}
				holds = append(holds, lock)
			}
		}
	}

	return holds
}

// AwaitingNotification ...
func AwaitingNotification(threads *[]ThreadInfo) map[Locked][]ThreadInfo {
	threadsWaiting := make(map[Locked][]ThreadInfo)

	for _, th := range *threads {
		if len(th.StackTrace) == 0 {
			continue
		}

		for _, stackLine := range strings.Split(th.StackTrace, "\n") {
			if rgxp, _ := regexp.Compile(parkingOrWaitingRgx); rgxp.MatchString(stackLine) {
				threadsWaiting = extractThreadWaiting(threadsWaiting, rgxp, stackLine, &th)
			}
		}

	}
	return threadsWaiting
}

func uniqueStackTrace(threadStackTrace []string) []string {
	u := make([]string, 0, len(threadStackTrace))
	m := make(map[string]bool)

	for _, val := range threadStackTrace {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func MostUsedMethods(threads *[]ThreadInfo) map[string]int {

	mostUsedMethodsGeneral := make(map[string]int)
	for _, th := range *threads {
		stackTraceLines := uniqueStackTrace(strings.Split(th.StackTrace, "\n"))
		for _, stackTraceLine := range stackTraceLines {
			if _, ok := mostUsedMethodsGeneral[stackTraceLine]; ok {
				mostUsedMethodsGeneral[stackTraceLine] = mostUsedMethodsGeneral[stackTraceLine] + 1
			} else {
				mostUsedMethodsGeneral[stackTraceLine] = 1
			}
		}

	}

	return mostUsedMethodsGeneral
}

func extractThreadWaiting(
	threadsWaiting map[Locked][]ThreadInfo,
	rgxp *regexp.Regexp, stackLine string,
	th *ThreadInfo) map[Locked][]ThreadInfo {

	for _, group := range rgxp.FindAllStringSubmatch(stackLine, -1) {
		k := Locked{group[1], group[2]}
		if _, exists := threadsWaiting[k]; !exists {
			threadsWaiting[k] = make([]ThreadInfo, 0)
		}
		h := threadsWaiting[k]
		h = append(h, *th)
		threadsWaiting[k] = h
	}
	return threadsWaiting

}
