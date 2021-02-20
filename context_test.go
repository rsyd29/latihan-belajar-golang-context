package latihan_belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
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

// Context With Cancel
//func sebelum menggunakan parameter
func CreateCounter() chan int { // membuat sebuah func CreateCounter dengan return value chan int
	destination := make(chan int) // membuat sebuah channel dengan tipe data int dan dimasukkan ke dalam var destination

	go func() { // membuat goroutine dan menjalankan sebuah anonymous function
		defer close(destination) // melakukan defer untuk close sebuah channel apabila proses goroutine sudah selesai
		counter := 1             // membuat konstanta dengan var counter yang bernilai awal 1
		for {                    // melakukan perulangan
			destination <- counter // memasukkan data counter ke dalam channel destination
			counter++              // increment sebuah counter yang nilainya akan terus bertambah
		}
	}() // menjalankan sebuah anonymous function

	return destination // ini akan mengembalikan nilai destination
}

// func setelah menambahkan parameter
func CreateCounterWithParameter(ctx context.Context) chan int { // membuat sebuah func CreateCounter dengan return value chan int
	destination := make(chan int) // membuat sebuah channel dengan tipe data int dan dimasukkan ke dalam var destination

	go func() { // membuat goroutine dan menjalankan sebuah anonymous function
		defer close(destination) // melakukan defer untuk close sebuah channel apabila proses goroutine sudah selesai
		counter := 1             // membuat konstanta dengan var counter yang bernilai awal 1
		for {                    // melakukan perulangan
			select { // melakukan select
			case <-ctx.Done(): // apabila kasusnya context nya sudah selesai maka
				return // hentikan perulangan
			default: // apabila belum selesai
				destination <- counter // memasukkan data counter ke dalam channel destination
				counter++              // increment sebuah counter yang nilainya akan terus bertambah
			}

		}
	}() // menjalankan sebuah anonymous function

	return destination // ini akan mengembalikan nilai destination
}

func TestContextWithCancel(t *testing.T) { // membuat unit test TestContextWithCancel
	fmt.Println("Total Goroutine Awal", runtime.NumGoroutine()) // cetak goroutine awal sebelum melakukan perulangan

	destination := CreateCounter() // membuat sebuah var destination dengan nilai sebuah function CreateCounter()

	fmt.Println("Total Goroutine Func", runtime.NumGoroutine()) // cetak goroutine awal sebelum melakukan perulangan

	for n := range destination { // melakukan for range perulangan dengan rangenya itu ada di destination yang merupakan sebuah counter yang terus bertambah
		fmt.Println("Counter", n) // cetak sebuah counter yang berisikan nilainya berupa angka
		if n == 10 {              // jika nilai n = 10, maka
			break // program akan berhenti, dan kita tidak butuh lagi goroutine
		}
	}
	time.Sleep(2 * time.Second) // untuk menunggu goroutine-nya selesai selama 2 detik

	fmt.Println("Total Goroutine Akhir", runtime.NumGoroutine()) // selanjutnya menghitung total goroutine akhir, setelah program berhasil dijalankan

	/**
	Outputnya terjadi goroutine leak, jadi goroutinenya jalan terus padahal program kita tidak membutuhkan goroutine-nya
	=== RUN   TestContextWithCancel
	Total Goroutine Awal 2
	Total Goroutine Func 3
	Counter 1
	Counter 2
	Counter 3
	Counter 4
	Counter 5
	Counter 6
	Counter 7
	Counter 8
	Counter 9
	Counter 10
	Total Goroutine Akhir 3
	--- PASS: TestContextWithCancel (2.00s)
	PASS

	Seharusnya total goroutine akhir itu kembali ke 2 karena sudah tidak dibutuhkan lagi goroutine karena ada proses break
	Jadi 1 goroutine itu ada yang menyala terus tidak pernah mati, kenapa?
	karena pada func CreateCounter() itu membuat sebuah perulangan yang tidak berhenti-henti, dan kita selalu mengirimkan
	data counter-nya ke dalam channel destination. Artinya apa?
	Ada yang consume channel destination ataupun tidak nah si goroutine ini tetap mencoba mengirimkan datanya ke dalam
	channel. Nah di func unit test kita hanya meng-consue dan selesai, artinya
	goroutine-nya terjadi leak, goroutine jalan terus tanpa berhenti-henti.

	Ini berbahaya apabila ada goroutine yang jalan terus menerus, biasanya apabila banyak goroutine yang leak maka
	program yang kita buat akan menjadi lambat dan memory consume-nya akan semakin tinggi lalu aplikasinya mati dan tidak
	sadar kenapa bisa mati, padahal itu karena terkena goroutine leak.

	Alangkah baiknya kita memakai fitur signal context with cancel.
	*/
}

func TestContextWithCancelParameter(t *testing.T) {
	fmt.Println("Total Goroutine Awal", runtime.NumGoroutine()) // cetak goroutine awal sebelum melakukan perulangan
	parent := context.Background()                              // membuat context background
	ctx, cancel := context.WithCancel(parent)                   // membuat context with cancel, mengembalikan 2 value yaitu ctx Context dan cancel function

	destinationWithParameter := CreateCounterWithParameter(ctx) // mengambil data dari func yang memiliki parameter ctx

	fmt.Println("Total Goroutine Func", runtime.NumGoroutine()) // cetak goroutine awal sebelum melakukan perulangan

	for n := range destinationWithParameter { // melakukan for range terhadap func dengan parameter
		fmt.Println("Counter", n) // cetak sebuah counter yang berisikan nilainya berupa angka
		if n == 10 {              // apabila n-nya 10
			break // maka break perulangan
		}
	}
	cancel() // mengirim sinyal cancel ke context

	time.Sleep(2 * time.Second) // untuk menunggu goroutine-nya selesai selama 2 detik

	fmt.Println("Total Goroutine Akhir", runtime.NumGoroutine()) // selanjutnya menghitung total goroutine akhir, setelah program berhasil dijalankan

	/**
	Outputnya setelah menggunakan context cancel
	=== RUN   TestContextWithCancelParameter
	Total Goroutine Awal 2
	Total Goroutine Func 3
	Counter 1
	Counter 2
	Counter 3
	Counter 4
	Counter 5
	Counter 6
	Counter 7
	Counter 8
	Counter 9
	Counter 10
	Total Goroutine Akhir 2
	--- PASS: TestContextWithCancelParameter (2.00s)
	PASS
	*/
}

// Context With Timeout
func CreateCounterWithTimeout(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulasi slow
			}
		}
	}()
	return destination
}

// Unit test untuk context with timeout
func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine Awal", runtime.NumGoroutine())
	parent := context.Background()
	// context.WithTimeout(parent, duration)
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel() // untuk memastikan function unit test sudah dijalankan maka cancel akan dieksekusi

	destination := CreateCounterWithTimeout(ctx)
	fmt.Println("Total Goroutine Func", runtime.NumGoroutine())
	for n := range destination { // melakukan looping tidak pernah berhenti, jadi apabila itu tidak pernah berhenti
		// sampai waktu duration timeout blm seleesai juga, maka akan dibatalkan
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine Akhir", runtime.NumGoroutine())
	/**
	Output Program
	=== RUN   TestContextWithTimeout
	Total Goroutine Awal 2
	Total Goroutine Func 3
	Counter 1
	Counter 2
	Counter 3
	Counter 4
	Counter 5
	Total Goroutine Akhir 2
	--- PASS: TestContextWithTimeout (7.00s)
	PASS

	Kenapa keluar counternya hanya sampai 5? karena kita berhentinya setiap 1 detik (slow motion),
	jadi setiap counter akan sleep selama 1 detik sampai counter 5. Ketika sudah sampai counter 5.
	Duration timeoutnya hanya sampai 5 detik, maka pada saat counter 5 sinyal cancel dikirim masuk ke dalam select yang
	menyatakan bahwa ctx.Done() telah selesai dan dibatalkan, dan channel destination akan di close.

	Jadi dengan menggunakan duration timeout, maka program counter itu tidak akan lebih dari 5 detik.
	Jadi tetap direkomendasikan tetap memanggil
	defer cancel()
	*/
}
