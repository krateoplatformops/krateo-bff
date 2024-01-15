package evaluator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/krateoplatformops/krateo-bff/internal/tmpl"
)

func TestIteratorCount(t *testing.T) {
	sample := `{"items": [
	{
		"postId": 1,
		"id": 1,
		"name": "id labore ex et quam laborum",
		"email": "Eliseo@gardner.biz",
		"body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
	},
	{
		"postId": 1,
		"id": 2,
		"name": "quo vero reiciendis velit similique earum",
		"email": "Jayne_Kuhic@sydney.com",
		"body": "est natus enim nihil est dolore omnis voluptatem numquam\net omnis occaecati quod ullam at\nvoluptatem error expedita pariatur\nnihil sint nostrum voluptatem reiciendis et"
	},
	{
		"postId": 1,
		"id": 3,
		"name": "odio adipisci rerum aut animi",
		"email": "Nikita@garfield.biz",
		"body": "quia molestiae reprehenderit quasi aspernatur\naut expedita occaecati aliquam eveniet laudantium\nomnis quibusdam delectus saepe quia accusamus maiores nam est\ncum et ducimus et vero voluptates excepturi deleniti ratione"
	}
]}
`

	ds := map[string]any{}
	err := json.Unmarshal([]byte(sample), &ds)
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		t.Fatal(err)
	}

	tot := 1

	it := ".items"
	if len(it) > 0 {
		len, err := tpl.Execute(fmt.Sprintf("${ %s | length }", it), ds)
		if err != nil {
			t.Fatal(err)
		}
		tot, err = strconv.Atoi(len)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("tot = ", tot)
}

func TestIterator(t *testing.T) {
	sample := `{"items": [
	{
		"postId": 1,
		"id": 1,
		"name": "id labore ex et quam laborum",
		"email": "Eliseo@gardner.biz",
		"body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
	},
	{
		"postId": 1,
		"id": 2,
		"name": "quo vero reiciendis velit similique earum",
		"email": "Jayne_Kuhic@sydney.com",
		"body": "est natus enim nihil est dolore omnis voluptatem numquam\net omnis occaecati quod ullam at\nvoluptatem error expedita pariatur\nnihil sint nostrum voluptatem reiciendis et"
	},
	{
		"postId": 1,
		"id": 3,
		"name": "odio adipisci rerum aut animi",
		"email": "Nikita@garfield.biz",
		"body": "quia molestiae reprehenderit quasi aspernatur\naut expedita occaecati aliquam eveniet laudantium\nomnis quibusdam delectus saepe quia accusamus maiores nam est\ncum et ducimus et vero voluptates excepturi deleniti ratione"
	}
]}
`

	ds := map[string]any{}
	err := json.Unmarshal([]byte(sample), &ds)
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := tmpl.New("${", "}")
	if err != nil {
		t.Fatal(err)
	}

	tot := 1

	it := ".items"
	if len(it) > 0 {
		len, err := tpl.Execute(fmt.Sprintf("${ %s | length }", it), ds)
		if err != nil {
			t.Fatal(err)
		}
		tot, err = strconv.Atoi(len)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("tot = ", tot)
}
