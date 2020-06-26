package utfutil_test

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/TomOnTime/utfutil"
)

func TestReadFile(t *testing.T) {

	expected, err := ioutil.ReadFile(filepath.Join("testdata", "calblur8.htm"))
	if err != nil {
		log.Fatal(err)
	}

	// The test files were generated with:
	//    for i in $(iconv  -l|grep UTF)  ; do
	//        iconv -f UTF-8 -t $i calblur8.htm > calblur8.htm.$i
	//    done
	for _, tst := range []struct {
		works bool // is combination is expected to work?
		ume   utfutil.EncodingHint
		name  string
	}{
		// Assume missing BOM means UTF8
		{true, utfutil.UTF8, "calblur8.htm.UTF-8"},     // No BOM
		{true, utfutil.UTF8, "calblur8.htm.UTF-16"},    // BOM=fffe
		{false, utfutil.UTF8, "calblur8.htm.UTF-16LE"}, // no BOM
		{false, utfutil.UTF8, "calblur8.htm.UTF-16BE"}, // no BOM
		// Assume missing BOM means UFT16LE
		{false, utfutil.UTF16LE, "calblur8.htm.UTF-8"},    // No BOM
		{true, utfutil.UTF16LE, "calblur8.htm.UTF-16"},    // BOM=fffe
		{true, utfutil.UTF16LE, "calblur8.htm.UTF-16LE"},  // no BOM
		{false, utfutil.UTF16LE, "calblur8.htm.UTF-16BE"}, // no BOM
		// Assume missing BOM means UFT16BE
		{false, utfutil.UTF16BE, "calblur8.htm.UTF-8"},    // No BOM
		{true, utfutil.UTF16BE, "calblur8.htm.UTF-16"},    // BOM=fffe
		{false, utfutil.UTF16BE, "calblur8.htm.UTF-16LE"}, // no BOM
		{true, utfutil.UTF16BE, "calblur8.htm.UTF-16BE"},  // no BOM
	} {

		actual, err := utfutil.ReadFile(filepath.Join("testdata", tst.name), tst.ume)
		if err != nil {
			log.Fatal(err)
		}

		if tst.works {
			if string(expected) == string(actual) {
				t.Log("SUCCESS:", tst.ume, tst.name)
			} else {
				t.Errorf("FAIL: %v/%v: expected %#v got %#v\n", tst.ume, tst.name, string(expected)[:4], actual[:4])
			}
		} else {
			if string(expected) != string(actual) {
				t.Logf("SUCCESS: %v/%v: failed as expected.", tst.ume, tst.name)
			} else {
				t.Errorf("FAILUREish: %v/%v: unexpected success!", tst.ume, tst.name)
			}
		}
	}

}

func TestReadAndCloseFile(t *testing.T) {
	file := filepath.Join("testdata", "calblur8.htm.UTF-8")
	_, err := utfutil.ReadFile(file, utfutil.UTF8)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_EXCL, 0)
	if err != nil {
		t.Errorf("FAIL: Unable to open file in exclusive mode after reading, handle must still be open\n")
	}

	f.Close()
	t.Log("SUCCESS: Closed file after reading")
}

func TestReadAndCloseScanner(t *testing.T) {
	file := filepath.Join("testdata", "calblur8.htm.UTF-8")
	scanner, err := utfutil.NewScanner(file, utfutil.UTF8)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		scanner.Text()
	}

	if err := scanner.Close(); err != nil {
		t.Errorf("FAIL: Unable to close file handle after scan")
	}

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_EXCL, 0)
	if err != nil {
		t.Errorf("FAIL: Unable to open file in exclusive mode after reading, handle must still be open")
	}

	if err := f.Close(); err != nil {
		t.Errorf("FAIL: Unable to close file handle after reading")
	}

	t.Logf("SUCCESS: Read and closed file handle")
}

func TestReadAndCloseFileReader(t *testing.T) {
	file := filepath.Join("testdata", "calblur8.htm.UTF-8")
	fr, err := utfutil.OpenFile(file, utfutil.UTF8)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(fr)
	for {
		_, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}
	}

	if err := fr.Close(); err != nil {
		t.Errorf("FAIL: Unable to close file handle after reading")
	}

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_EXCL, 0)
	if err != nil {
		t.Errorf("FAIL: Unable to open file in exclusive mode after reading, handle must still be open\n")
	}

	if err := f.Close(); err != nil {
		t.Errorf("FAIL: Unable to close file handle after reading")
	}

	t.Logf("SUCCESS: Read and closed file handle")
}
