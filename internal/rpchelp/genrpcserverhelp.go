// Copyright (c) 2015 Conformal Systems LLC <info@conformal.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

//+build generate

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcjson/v2/btcjson"
	"github.com/btcsuite/btcwallet/internal/rpchelp"
)

var outputFile = func() *os.File {
	fi, err := os.Create("rpcserverhelp.go")
	if err != nil {
		log.Fatal(err)
	}
	return fi
}()

func writefln(format string, args ...interface{}) {
	_, err := fmt.Fprintf(outputFile, format, args...)
	if err != nil {
		log.Fatal(err)
	}
	_, err = outputFile.Write([]byte{'\n'})
	if err != nil {
		log.Fatal(err)
	}
}

func writeLocaleHelp(locale, goLocale string, descs map[string]string) {
	funcName := "helpDescs" + goLocale
	writefln("func %s() map[string]string {", funcName)
	writefln("return map[string]string{")
	for i := range rpchelp.Methods {
		m := &rpchelp.Methods[i]
		helpText, err := btcjson.GenerateHelp(m.Method, descs, m.ResultTypes...)
		if err != nil {
			log.Fatal(err)
		}
		writefln("%q: %q,", m.Method, helpText)
	}
	writefln("}")
	writefln("}")
}

func writeLocales() {
	writefln("var localeHelpDescs = map[string]func() map[string]string{")
	for _, h := range rpchelp.HelpDescs {
		writefln("%q: helpDescs%s,", h.Locale, h.GoLocale)
	}
	writefln("}")
}

func writeUsage() {
	usageStrs := make([]string, len(rpchelp.Methods))
	var err error
	for i := range rpchelp.Methods {
		usageStrs[i], err = btcjson.MethodUsageText(rpchelp.Methods[i].Method)
		if err != nil {
			log.Fatal(err)
		}
	}
	usages := strings.Join(usageStrs, "\n")
	writefln("var requestUsages = %q", usages)
}

func main() {
	defer outputFile.Close()

	writefln("// AUTOGENERATED by internal/rpchelp/genrpcserverhelp.go; do not edit.")
	writefln("")
	writefln("package main")
	writefln("")
	for _, h := range rpchelp.HelpDescs {
		writeLocaleHelp(h.Locale, h.GoLocale, h.Descs)
		writefln("")
	}
	writeLocales()
	writefln("")
	writeUsage()
}
