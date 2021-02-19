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

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	// context.WithValue(parent, key, value)

	// parent-nya contextA
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	// parent-nya contextB
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	// parent-nya contextC
	contextF := context.WithValue(contextC, "f", "F")

	// parent-nya contextF
	contextG := context.WithValue(contextF, "g", "G")

	// cetak data context
	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	/**
	Output program
	=== RUN   TestContextWithValue
	context.Background
	context.Background.WithValue(type string, val B)
	context.Background.WithValue(type string, val C)
	context.Background.WithValue(type string, val B).WithValue(type string, val D)
	context.Background.WithValue(type string, val B).WithValue(type string, val E)
	context.Background.WithValue(type string, val C).WithValue(type string, val F)
	context.Background.WithValue(type string, val C).WithValue(type string, val F).WithValue(type string, val G)
	--- PASS: TestContextWithValue (0.00s)
	PASS
	*/

	/**
	Cara pengambilan value di context, yaitu dengan cara melihat context child terlebih dahulu apabila tidak ada maka
	akan melihat parent-nya yang paling tinggi. Kalau memang tidak ada berarti nilainya adalah nil.
	*/
	fmt.Println("\nContext Get Value")
	fmt.Println(contextF.Value("f")) // dapat
	fmt.Println(contextF.Value("c")) // dapat milik parent
	fmt.Println(contextF.Value("b")) // tidak dapat, beda parent
	fmt.Println(contextA.Value("b")) // tidak bisa mengambil data child, kenapa nil karena context itu menanyakannya ke atas, maka ketika dia tidak ada maka menanyakan ke parent, dia tidak mungkin menanyakan ke child-nya
	/**
	Output program
	Context Get Value
	F
	C
	<nil>
	<nil>

	Jadi value itu pasti menanyakannya itu ke parent-nya bukan ke child-nya
	*/
}
