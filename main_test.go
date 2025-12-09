package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		main()

		return
	}

	os.Exit(m.Run())
}

func runCPUTest(t *testing.T, testFileName, expected string) {
	cmd := exec.Command(os.Args[0], "--cpm", "--unthrottle", "--headless", "./sub/8080/cpu_tests/"+testFileName)

	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	out, err := cmd.Output()
	assert.NoError(t, err)

	t.Logf("Output from %s:\n%s", testFileName, string(out))
	assert.Equal(t, expected, string(out))
}

func Test_CPU(t *testing.T) {
	t.Run("TST8080.COM", func(t *testing.T) {
		expected := "MICROCOSM ASSOCIATES 8080/8085 CPU DIAGNOSTIC\r\n VERSION 1.0  (C) 1980\r\n\r\n CPU IS OPERATIONAL"
		runCPUTest(t, "TST8080.COM", expected)
	})

	t.Run("8080PRE.COM", func(t *testing.T) {
		expected := "8080 Preliminary tests complete"
		runCPUTest(t, "8080PRE.COM", expected)
	})

	t.Run("CPUTEST.COM", func(t *testing.T) {
		expected := "\x00\x00\x00\x00\x00\x00\r\nDIAGNOSTICS II V1.2 - CPU TEST\r\nCOPYRIGHT (C) 1981 - SUPERSOFT ASSOCIATES\r\n\nABCDEFGHIJKLMNOPQRSTUVWXYZ\r\nCPU IS 8080/8085\r\nBEGIN TIMING TEST\r\n\a\aEND TIMING TEST\r\nCPU TESTS OK\r\n"
		runCPUTest(t, "CPUTEST.COM", expected)
	})

	t.Run("8080EXM.COM", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		expected := "8080 instruction exerciser\n\rdad <b,d,h,sp>................  PASS! crc is:14474ba6\n\raluop nn......................  PASS! crc is:9e922f9e\n\raluop <b,c,d,e,h,l,m,a>.......  PASS! crc is:cf762c86\n\r<daa,cma,stc,cmc>.............  PASS! crc is:bb3f030c\n\r<inr,dcr> a...................  PASS! crc is:adb6460e\n\r<inr,dcr> b...................  PASS! crc is:83ed1345\n\r<inx,dcx> b...................  PASS! crc is:f79287cd\n\r<inr,dcr> c...................  PASS! crc is:e5f6721b\n\r<inr,dcr> d...................  PASS! crc is:15b5579a\n\r<inx,dcx> d...................  PASS! crc is:7f4e2501\n\r<inr,dcr> e...................  PASS! crc is:cf2ab396\n\r<inr,dcr> h...................  PASS! crc is:12b2952c\n\r<inx,dcx> h...................  PASS! crc is:9f2b23c0\n\r<inr,dcr> l...................  PASS! crc is:ff57d356\n\r<inr,dcr> m...................  PASS! crc is:92e963bd\n\r<inx,dcx> sp..................  PASS! crc is:d5702fab\n\rlhld nnnn.....................  PASS! crc is:a9c3d5cb\n\rshld nnnn.....................  PASS! crc is:e8864f26\n\rlxi <b,d,h,sp>,nnnn...........  PASS! crc is:fcf46e12\n\rldax <b,d>....................  PASS! crc is:2b821d5f\n\rmvi <b,c,d,e,h,l,m,a>,nn......  PASS! crc is:eaa72044\n\rmov <bcdehla>,<bcdehla>.......  PASS! crc is:10b58cee\n\rsta nnnn / lda nnnn...........  PASS! crc is:ed57af72\n\r<rlc,rrc,ral,rar>.............  PASS! crc is:e0d89235\n\rstax <b,d>....................  PASS! crc is:2b0471e9\n\rTests complete"
		runCPUTest(t, "8080EXM.COM", expected)
	})
}
