package batRun_test

import (
	"errors"
	"github.com/uynap/batRun"
	"testing"
	"time"
)

func TestBatRunWithoutTimeout(t *testing.T) {
	bat := batRun.NewBat()
	bat.AddProducers(func(task *batRun.Task) {
		for i := 1; i <= 2; i++ {
			data := map[string]int{
				"id": i,
			}
			ctx := task.NewContext(data)
			task.Submit(ctx, 0)
		}
	}, func(task *batRun.Task) {
		for i := 1; i <= 2; i++ {
			data := map[string]int{
				"id": i * 10,
			}
			ctx := task.NewContext(data)
			task.Submit(ctx, 0)
		}
	})

	bat.AddWork(func(ctx *batRun.Context) error {
		data := ctx.GetContext().(map[string]int)
		data["age"] = 22
		ctx.SetContext(data)
		return nil
	}, 2)

	bat.AddWork(func(ctx *batRun.Context) error {
		data := ctx.GetContext().(map[string]int)
		data["height"] = 182
		ctx.SetContext(data)
		return nil
	}, 4)

	bat.AddWork(func(ctx *batRun.Context) error {
		data := ctx.GetContext().(map[string]int)
		if data["age"] != 22 || data["height"] != 182 {
			t.Error("context error")
		}
		ctx.SetContext(data)
		return nil
	}, 1)

	bat.Run()
}

func TestBatRunWithTimeout(t *testing.T) {
	bat := batRun.NewBat()
	bat.AddProducers(func(task *batRun.Task) {
		for i := 1; i <= 2; i++ {
			data := map[string]int{
				"id": i,
			}
			ctx := task.NewContext(data)
			task.Submit(ctx, 3*time.Second) // create a task with timeout in 3s
		}
	})

	bat.AddWork(func(ctx *batRun.Context) error {
		data := ctx.GetContext().(map[string]int)
		data["age"] = 22
		ctx.SetContext(data)
		time.Sleep(1 * time.Second)
		return nil
	}, 2)

	bat.AddWork(func(ctx *batRun.Context) error {
		data := ctx.GetContext().(map[string]int)
		data["height"] = 182
		ctx.SetContext(data)
		time.Sleep(3 * time.Second)
		return nil
	}, 4)

	bat.Run()
}

func TestBatRunWithTimeoutCancel(t *testing.T) {
	bat := batRun.NewBat()
	bat.AddProducers(func(task *batRun.Task) {
		for i := 1; i <= 2; i++ {
			data := map[string]int{
				"id": i,
			}
			ctx := task.NewContext(data)
			task.Submit(ctx, 3*time.Second) // create a task with timeout in 3s
		}
	})

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 1")
		}
		data := ctx.GetContext().(map[string]int)
		data["age"] = 22
		ctx.SetContext(data)
		time.Sleep(1 * time.Second)
		return nil
	}, 5)

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 2")
		}
		data := ctx.GetContext().(map[string]int)
		data["size"] = 40
		ctx.SetContext(data)
		time.Sleep(5 * time.Second)
		return nil
	}, 5)

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 3")
		}
		data := ctx.GetContext().(map[string]int)
		data["height"] = 182
		ctx.SetContext(data)
		time.Sleep(3 * time.Second)
		return nil
	}, 5)

	bat.Run()
}

func TestBatRunWithCancel(t *testing.T) {
	bat := batRun.NewBat()
	bat.AddProducers(func(task *batRun.Task) {
		ctx := task.NewContext(map[string]int{})
		task.Submit(ctx, 0)
	})

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 1")
		}
		return nil
	}, 5)

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 2")
		}
		return errors.New("with some errors")
	}, 5)

	bat.AddWork(func(ctx *batRun.Context) error {
		ctx.Cancel = func() {
			println("Cancel is called from work 3")
		}
		return nil
	}, 5)

	bat.Run()
}
