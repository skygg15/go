package main

import (
 "fmt"
 "time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
 id         int
 cT         string // время создания
 fT         string // время выполнения
 taskRESULT []byte
}

func main() {
 taskCreator := func(a chan Ttype) {
  go func() {
   for {
    ft := time.Now().Format(time.RFC3339)
    if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
     ft = "Some error occurred"
    }
    a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
   }
  }()
 }

 superChan := make(chan Ttype, 10)

 go taskCreator(superChan)

 taskWorker := func(a Ttype) Ttype {
  tt, _ := time.Parse(time.RFC3339, a.cT)
  if tt.After(time.Now().Add(-20 * time.Second)) {
   a.taskRESULT = []byte("task has been succeeded")
  } else {
   a.taskRESULT = []byte("something went wrong")
  }
  a.fT = time.Now().Format(time.RFC3339Nano)

  time.Sleep(time.Millisecond * 150)

  return a
 }

 doneTasks := make(chan Ttype)
 undoneTasks := make(chan error)

 taskSorter := func(t Ttype) {
  if string(t.taskRESULT[14:]) == "succeeded" {
   doneTasks <- t
  } else {
   undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
  }
 }

 go func() {
  // получение тасков
  for t := range superChan {
   t = taskWorker(t)
   go taskSorter(t)
  }
  close(superChan)
 }()

 result := map[int]Ttype{}
 err := []error{}
 go func() {
  for r := range doneTasks {
   go func(r Ttype) {
    result[r.id] = r
   }(r)
  }
  for r := range undoneTasks {
   go func(r error) {
    err = append(err, r)
   }(r)
  }
  close(doneTasks)
  close(undoneTasks)
 }()

 time.Sleep(time.Second * 3)

 println("Errors:")
 for _, r := range err {
  println(r)
 }

 println("Done tasks:")
 for r := range result {
  println(r)
 }
}