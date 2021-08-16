package main

// func main() {
// 	c := sync.NewCond(&sync.Mutex{})

// 	var ready int
// 	for i := 0; i < 10; i++ {
// 		go func(i int) {
// 			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
// 			//加锁更改等待条件
// 			c.L.Lock()
// 			ready++
// 			c.L.Unlock()

// 			log.Printf("运动员#%d 已经准备就绪 \n", i)

// 			if ready == 10 {
// 				//广播唤醒所有的等待者
// 				c.Signal()
// 			}

// 		}(i)
// 	}

// 	// c.L.Lock()

// 	c.Wait()
// 	log.Println("裁判员被唤醒")

// 	c.L.Unlock()

// 	log.Println("所有运动员都准备就绪 比赛开始  3,2,1 ......")
// }
