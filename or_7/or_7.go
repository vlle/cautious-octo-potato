package main


import (
  "time"
  "fmt"
)

func or (channels ...<- chan interface{}) <- chan interface{} {
  or_ch := make(chan interface{})
  for  {
    for _, val := range channels {
      select {
      case x, ok := <-val:
        if ok {
          fmt.Printf("Value %d was read.\n", x)
        } else {
          fmt.Println("Channel closed!")
          go func() {
            fmt.Println("Channel closed!")
            defer close(or_ch)
          }()
          return or_ch
        }
      default:
        fmt.Println("No value ready, moving on.")
      }
    }
  }
  return or_ch
}

func main() {
  sig := func(after time.Duration) <- chan interface{} {
    c := make(chan interface{})
    go func() {
      defer close(c)
      time.Sleep(after)
    }()
    return c
  }

  start := time.Now()
  <-or (
    sig(2*time.Hour),
    sig(5*time.Minute),
    sig(4*time.Second),
    sig(1*time.Hour),
    sig(1*time.Minute),
  )

  fmt.Printf("fone after %v", time.Since(start))

}
