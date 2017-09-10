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
    threadNameRgxGroupIndex     = 1
    threadPriorityRgxGroupIndex = 2
    threadIdRgxGroupIndex       = 3
    threadNativeIdRgxGroupIndex = 4
)

type ThreadInfo struct {
    Name, ID, NativeID, Priority, State, StackTrace string
}

func (th ThreadInfo) String() string {
    return fmt.Sprintf("Thread Id: '%s', Name: '%s', NativeID: '%s', Priority: '%s', State: '%s', StackTrace: \n'%s'",
        th.ID, th.Name, th.NativeID, th.Priority, th.Priority, th.State, th.StackTrace)
}

func extractThreadInfoFromLine(line string) ThreadInfo {
    ti := ThreadInfo{}
    if rgxp, _ := regexp.Compile(threadNameRgx); rgxp.MatchString(line) {
        for _, v := range rgxp.FindAllStringSubmatch(line, -1) {
            ti.Name, ti.Priority, ti.ID, ti.NativeID = 
                v[threadNameRgxGroupIndex], v[threadPriorityRgxGroupIndex], v[threadIdRgxGroupIndex], v[threadNativeIdRgxGroupIndex]
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
