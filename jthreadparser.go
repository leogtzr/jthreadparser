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
	threadHeaderInfoRgx         = `^\"(.*).*tid=(\w*) nid=(\w*)\s\w*`
	threadNameRgx               = `^"(.*)".*\stid=(\w*) nid=(\w*)\s\w*`
	threadNameWithPriorityRgx   = `^"(.*)".*\sprio=(\d+).*\stid=(\w*) nid=(\w*)\s\w*`
	stateRgx                    = `\s*java.lang.Thread.State: (.*)`
	lockedRgx                   = `\s*\- locked\s*<(.*)>\s*\(a\s(.*)\)`
	runnableStateRgx            = `runnable\s{1,2}$`
	waitingOnStateRgx           = `waiting on condition\s{1,2}$`
	parkingOrWaitingRgx         = `\s*\- (?:waiting on|parking to wait for)\s*<(.*)>\s*\(a\s(.*)\)`
	stackTraceRgx               = `^\s+(at|\-\s).*\)$`
	stackTraceRgxMethodName     = `at\s+(.*)$`
	threadNameRgxGroupIndex     = 1
	threadPriorityRgxGroupIndex = 2
	threadIDRgxGroupIndex       = 3
	threadNativeIDRgxGroupIndex = 4
)

func (th ThreadInfo) String() string {
	if th.Daemon {
		return fmt.Sprintf("Thread Id: '%s' (daemon), Name: '%s', State: '%s'", th.ID, th.Name, th.State)
	}
	return fmt.Sprintf("Thread Id: '%s', Name: '%s', State: '%s'", th.ID, th.Name, th.State)
}

func extractThreadInfoFromLine(line string) ThreadInfo {
	ti := ThreadInfo{}

	ti.Daemon = strings.Contains(line, " daemon ")
	if strings.Contains(line, " prio=") {
		if rgxp, _ := regexp.Compile(threadNameWithPriorityRgx); rgxp.MatchString(line) {
			for _, v := range rgxp.FindAllStringSubmatch(line, -1) {
				ti.Name, ti.Priority, ti.ID, ti.NativeID =
					v[threadNameRgxGroupIndex], v[threadPriorityRgxGroupIndex], v[threadIDRgxGroupIndex], v[threadNativeIDRgxGroupIndex]
			}
		}
	} else {
		if rgxp, _ := regexp.Compile(threadNameRgx); rgxp.MatchString(line) {
			for _, v := range rgxp.FindAllStringSubmatch(line, -1) {
				ti.Name, ti.ID, ti.NativeID = v[threadNameRgxGroupIndex], v[2], v[3]
			}
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
	parse2(file, &threads)

	return threads, nil
}

func ParseFromFile2(fileName string) ([]ThreadInfo, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	threads := make([]ThreadInfo, 0)
	parse2(file, &threads)

	return threads, nil
}

// ParseFromFile ...
// func ParseFromFile2(fileName string) ([]ThreadInfo, error) {

// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	threads := make([]ThreadInfo, 0)
// 	parse(file, &threads)

// 	return threads, nil
// }

func parseOLD(r io.Reader, threads *[]ThreadInfo) {
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

func hasRunnableState(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(runnableStateRgx)
	return rgxp.MatchString(threadHeaderLine)
}

func parse2(r io.Reader, threads *[]ThreadInfo) {
	scanner := bufio.NewScanner(r)
	tlines := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		tlines = append(tlines, line)
	}

	for i := 0; i < len(tlines); i++ {
		line := tlines[i]
		if hasThreadHeaderInformation(line) {
			threadInfo := extractThreadInfoFromLine(line)

			if hasRunnableState(line) {
				threadInfo.State = "runnable"
			}

			i++
			if i < len(tlines) {
				line = tlines[i]
			} else {
				break
			}
			// Look for the thread state:
			if threadState := extractThreadState(line); len(threadState) != 0 {
				threadInfo.State = threadState
				if i < len(tlines) {
					line = tlines[i]
					i++
				}
			} else {
				i--
			}

			// There could be two threads together without a thread state ...
			if hasThreadHeaderInformation(line) {
				line2 := line
				threadInfo2 := extractThreadInfoFromLine(line2)
				if hasRunnableState(line2) {
					threadInfo2.State = "runnable"
				} else if hasWaitingOnState(line2) {
					threadInfo2.State = "waiting on condition"
				}
				*threads = append(*threads, threadInfo2)
				i++
			}

			if i < len(tlines) {
				line = tlines[i]
			} else {
				break
			}

			var sb strings.Builder
			for (len(line) > 0) && !strings.HasPrefix(line, threadInformationBegins) {
				sb.WriteString(strings.TrimSpace(line))
				sb.WriteString("\n")
				i++

				if i < len(tlines) {
					line = tlines[i]
				} else {
					break
				}
			}

			stackTrace := sb.String()
			if len(stackTrace) > 0 {
				threadInfo.StackTrace = stackTrace
			}

			*threads = append(*threads, threadInfo)
		}
	}
}

// ParseFrom ...
func ParseFrom(r io.Reader) ([]ThreadInfo, error) {
	threads := make([]ThreadInfo, 0)
	parse2(r, &threads)
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

func extractJavaMethodNameFromStackTraceLine(stacktraceLine string) string {
	if rgxp, _ := regexp.Compile(stackTraceRgxMethodName); rgxp.MatchString(stacktraceLine) {
		fields := strings.Fields(strings.TrimSpace(stacktraceLine))
		return strings.Join(fields[1:], " ")
	}
	return ""
}

// MostUsedMethods ...
func MostUsedMethods(threads *[]ThreadInfo) map[string]int {

	mostUsedMethodsGeneral := make(map[string]int)
	for _, th := range *threads {
		stackTraceLines := uniqueStackTrace(strings.Split(th.StackTrace, "\n"))
		for _, stackTraceLine := range stackTraceLines {
			stackTraceLine = extractJavaMethodNameFromStackTraceLine(stackTraceLine)
			if stackTraceLine == "" {
				continue
			}
			if _, ok := mostUsedMethodsGeneral[stackTraceLine]; ok {
				mostUsedMethodsGeneral[stackTraceLine]++
			} else {
				mostUsedMethodsGeneral[stackTraceLine] = 1
			}
		}

	}

	return mostUsedMethodsGeneral
}

// IdenticalStackTrace ...
func IdenticalStackTrace(threads *[]ThreadInfo) map[string]int {
	indenticalStackTrace := make(map[string]int)
	for _, th := range *threads {
		thStack := strings.TrimSpace(th.StackTrace)
		if _, ok := indenticalStackTrace[thStack]; ok {
			indenticalStackTrace[thStack]++
		} else {
			indenticalStackTrace[thStack] = 1
		}
	}
	return indenticalStackTrace
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

func hasWaitingOnState(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(waitingOnStateRgx)
	return rgxp.MatchString(threadHeaderLine)
}

func hasThreadHeaderInformation(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(threadHeaderInfoRgx)
	return rgxp.MatchString(threadHeaderLine)
}
