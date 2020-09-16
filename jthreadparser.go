package jthreadparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
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
		if st := strings.Split(line, ":"); st != nil {
			state := strings.TrimSpace(st[1])
			return strings.Split(state, " ")[0]
		}
	}
	return ""
}

// ParseFromFile ...
func ParseFromFile(fileName string) ([]ThreadInfo, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	threads := make([]ThreadInfo, 0)
	parse(file, &threads)

	return threads, nil
}

func hasRunnableState(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(runnableStateRgx)
	return rgxp.MatchString(threadHeaderLine)
}

func parse(r io.Reader, threads *[]ThreadInfo) {

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
			threadInfoTogether := ThreadInfo{}
			threadTogether := false

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
				threadInfoTogether = extractThreadInfoFromLine(line2)
				threadTogether = true
				if hasRunnableState(line2) {
					threadInfoTogether.State = "runnable"
				} else if hasWaitingOnState(line2) {
					threadInfoTogether.State = "waiting on condition"
				}
				// *threads = append(*threads, threadInfo2)
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

			if threadTogether {
				*threads = append(*threads, threadInfo)
				*threads = append(*threads, threadInfoTogether)
			} else {
				*threads = append(*threads, threadInfo)
			}
		}
	}
}

// ParseFrom ...
func ParseFrom(r io.Reader) ([]ThreadInfo, error) {
	threads := make([]ThreadInfo, 0)
	parse(r, &threads)
	return threads, nil
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

func hasWaitingOnState(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(waitingOnStateRgx)
	return rgxp.MatchString(threadHeaderLine)
}

func hasThreadHeaderInformation(threadHeaderLine string) bool {
	rgxp := regexp.MustCompile(threadHeaderInfoRgx)
	return rgxp.MatchString(threadHeaderLine)
}

func extractSynchronizers(stacktrace string) []Synchronizer {
	syncs := make([]Synchronizer, 0)
	rgxp := regexp.MustCompile(synchronizerRgx)

	for _, stackLine := range strings.Split(stacktrace, "\n") {
		if rgxp.MatchString(stackLine) {
			sync := Synchronizer{}
			for _, group := range rgxp.FindAllStringSubmatch(stackLine, -1) {
				sync.State = convertToSyncState(group[1])
				sync.ID = group[2]
				sync.ObjectName = group[3]
			}
			syncs = append(syncs, sync)
		}
	}

	return syncs
}

func equal(a, b []Synchronizer) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func convertToSyncState(textState string) SynchronizerState {

	if strings.Contains(textState, "waiting on") {
		return WaitingOnState
	} else if strings.Contains(textState, "locked") {
		return LockedState
	} else if strings.Contains(textState, "parking to wait for") {
		return ParkingToWaitForState
	} else if strings.Contains(textState, "waiting to lock") {
		return WaitingToLockState
	}

	return LockedState
}

// Synchronizers ....
func Synchronizers(threads *[]ThreadInfo) map[ThreadInfo][]Synchronizer {
	syncs := make(map[ThreadInfo][]Synchronizer)

	for _, th := range *threads {
		if len(th.StackTrace) == 0 {
			continue
		}

		threadSyncs := extractSynchronizers(th.StackTrace)
		if len(threadSyncs) > 0 {
			syncs[th] = threadSyncs
		}
	}

	return syncs
}
