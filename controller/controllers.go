package controller

import (
	"errors"
	"github.com/yanatan16/go-todo-app/model"
	"github.com/yanatan16/go-todo-app/model/todo"
	"net/url"
	"strconv"
)

// Add an item to the list by interpreting a url.Values
func add(app *model.TodoApp, params url.Values) (map[string]interface{}, error) {
	desc := params.Get("desc")
	if desc == "" {
		return nil, errors.New("Query parameter 'desc' required")
	}

	numReturn := make(chan int)
	req := todo.ItemRequest{
		Desc: desc,
		Num:  numReturn,
	}
	app.List.Add <- req

	// Wait for the return
	n := <-numReturn
	ret := make(map[string]interface{})
	ret["num"] = n
	return ret, nil
}

// Set an item to the done value interpreting a url.Values
func done(app *model.TodoApp, params url.Values) (map[string]interface{}, error) {
	nstr := params.Get("num")
	donestr := params.Get("done")

	if nstr == "" {
		return nil, errors.New("Query parameter 'num' required")
	} else if donestr == "" {
		return nil, errors.New("Query parameter 'done' required")
	}

	// Parse the int as base 10 and 32 bits
	n, err := strconv.ParseInt(nstr, 10, 32)
	if err != nil {
		return nil, errors.New("Query parameter 'num' must be an integer: " + err.Error())
	}

	// Parse the boolean
	done, err := strconv.ParseBool(donestr)
	if err != nil {
		return nil, errors.New("Query parameter 'done' must be a boolean: " + err.Error())
	}

	if int(n) >= len(app.List.Items) {
		return nil, errors.New("Query parameter 'num' out of range!")
	}

	// Perform the action
	it := app.List.Items[n]
	it.Done = done
	app.List.Set <- it

	// Return no error
	return nil, nil
}
