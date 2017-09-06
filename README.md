```go
threads, err := jthreadparser.Parse(os.Args[1])
if err != nil {
    log.Fatal(err)
}

for _, th := range threads {
    fmt.Println(th)
}
```