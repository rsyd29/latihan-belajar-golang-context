package latihan_belajar_golang_context

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	/**
	Membuat context kosong. Tidak pernah dibatalkan, tidak pernah timeout, dan tidak memiliki value apapun. Biasanya
	digunakan di main function atau dalam test, atau dalam awal proses request terjadi.
	*/
	background := context.Background() // diawal kita akan menggunakan ini dan membuatnya secara manual, kalau via web kita tidak perlu membuat manual.
	fmt.Println(background)

	/**
	Membuat context kosong seperti Background(), namun biasanya menggunakan ini ketika belum jelas context apa yang
	ingin digunakan.
	*/
	todo := context.TODO()
	fmt.Println(todo)
	/**
	Output-nya
	=== RUN   TestContext
	context.Background
	context.TODO
	--- PASS: TestContext (0.00s)
	PASS
	*/
}
