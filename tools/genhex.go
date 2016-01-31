package main

import (
	"flag"
	"os"
	"fmt"
	"io"
)

func writeHead(o io.Writer, pkgname string) error {
	_, err := fmt.Fprint(o, "package " + pkgname + "\n\n")
	return err
}

var kb = []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101, 102}

func writeBody(raw io.Reader, o io.Writer, numberEach int, vn string) error {
	_, err := fmt.Fprintf(o, "var %s []byte = []byte{\n", vn)
	if err != nil {
		return err
	}

	buff := make([]byte, numberEach)
	tmpv := []byte{0x30, 0x78, 0x00, 0x00, 0x2c, 0x20} //[0x**, ]

	for {
		n, err := raw.Read(buff)
		if err != nil && err != io.EOF {
			return err
		} else if err == io.EOF {
			break
		}
		for i := 0; i < n; i++ {
			tmpv[2] = kb[buff[i] >> 4]
			tmpv[3] = kb[buff[i] & 0x0f]
			_, err := o.Write(tmpv)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(o, "\n")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(o, "}\n")
	return err
}

func writeTail(o io.Writer) error {
	return nil
}

func write(i io.Reader, o io.Writer, pkgname string, vn string, numberEach int) error {
	if err := writeHead(o, pkgname); err != nil {
		return err
	}
	if err := writeBody(i, o, numberEach, vn); err != nil {
		return err
	}
	return writeTail(o);
}

func beginTask(ifn string, ofn string, pkgname string, vn string, numberEach int) error {
	ii, e := os.Open(ifn)
	if e != nil {
		return e
	}
	defer ii.Close()

	oo, e := os.OpenFile(ofn, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0664)
	if e != nil {
		return e
	}
	defer oo.Close()

	return write(ii, oo, pkgname, vn, numberEach)
}

func main() {
	var (
		ifn string
		numberInEachRow int
		packageName string
		varName string
	)

	flag.StringVar(&ifn, "if", "", "input file")
	flag.IntVar(&numberInEachRow, "n", 10, "hex number in each row, should bigger or equal 1")
	flag.StringVar(&packageName, "package", "main", "package name")
	flag.StringVar(&varName, "var", "hex", "var name")
	flag.Parse()

	if ifn == "" || numberInEachRow <= 0{
		if ifn == "" {
			fmt.Printf("argument \"if\" should not empty\n")
		}
		if numberInEachRow <= 0 {
			fmt.Printf("argument \"n\" should not less than 1\n")
		}
		flag.Usage()
		os.Exit(0)
	}

	if err := beginTask(ifn, ifn+".go", packageName, varName, numberInEachRow); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println("Output file is " + ifn + ".go")
}