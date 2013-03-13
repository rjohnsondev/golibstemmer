package stemmer

import "testing"
import "fmt"
import "runtime"
import "io/ioutil"
import "strings"


func TestInvalidLang(t *testing.T) {
    _, err := NewStemmer("englih")
    if err == nil {
        t.Errorf("Invalid algorithm accepted")
    }
}

func TestUtf8(t *testing.T) {
    stemmer, _ := NewStemmer("english")
    defer stemmer.Close()
    word := "♬"
    stemmed := stemmer.StemWord(word)
    if stemmed != word {
        t.Errorf("UTF-8 char: %s managled in stemming", word)
    }
}

func TestCMemoryLeak(t *testing.T) {
    contents, _ := ioutil.ReadFile("war of the worlds.txt")
    words := strings.Split(string(contents), " ")

    fmt.Printf("Testing for leaks\n"+
               "Memory will increase over time due to fragmentation, "+
               "however look out for consistently large increases\n")
    stats := &runtime.MemStats{ }

    for x := 0; x < 6; x++ {
        for _, word := range words {
            stemmer, _ := NewStemmer("english")
            _ = stemmer.StemWord(word)
            stemmer.Close()
        }
        runtime.GC()
        runtime.ReadMemStats(stats)
        fmt.Printf("Allocated memory: %d\n", stats.Alloc)
    }

}

func TestStem(t *testing.T) {
    stemmer, _ := NewStemmer("english")
    defer stemmer.Close()
    word := stemmer.StemWord("happy")
    if word != "happi" {
        t.Errorf("\"happy\" stemmed to %s, not \"happi\" as expected")
    }
}

func TestStemList(t *testing.T) {
    list := GetSupportedLanguages()
    found := false
    for _, lang := range list {
        if lang == "english" {
            found = true
            break
        }
    }
    if !found {
        t.Errorf("Unable to find stemmer \"english\" in supported language list!")
    }
}
