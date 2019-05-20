package jthreadparser

import (
	"bufio"
	"bytes"
	"fmt"
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
	threadNameRgxGroupIndex     = 1
	threadPriorityRgxGroupIndex = 2
	threadIDRgxGroupIndex       = 3
	threadNativeIDRgxGroupIndex = 4
)

// ThreadInfo ...
type ThreadInfo struct {
	Name, ID, NativeID, Priority, State, StackTrace string
}

// Locked ...
type Locked struct {
	LockID, LockecObjectName string
}

func (th ThreadInfo) String() string {
	return fmt.Sprintf("Thread Id: '%s', Name: '%s', State: '%s'", th.ID, th.Name, th.State)
}

func extractThreadInfoFromLine(line string) ThreadInfo {
	ti := ThreadInfo{}
	if rgxp, _ := regexp.Compile(threadNameRgx); rgxp.MatchString(line) {
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

func Parse(fileName string) ([]ThreadInfo, error) {

	file, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}

	threads := make([]ThreadInfo, 0, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), threadInformationBegins) {
			ti := extractThreadInfoFromLine(scanner.Text())
			scanner.Scan()
			ti.State = extractThreadState(scanner.Text())
			ti.StackTrace = extractThreadStackTrace(scanner)
			threads = append(threads, ti)
		}

	}

	return threads, nil

}

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
