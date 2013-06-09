package stemmer

//import "runtime"
import "unsafe"
import "strings"
import "fmt"
import "sync"

/*
#cgo LDFLAGS: -lstemmer

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <stdlib.h>
#include <libstemmer.h>

// make a copy of the passed in c string ready for
// use by the stemmer, it must be freed after use.
sb_symbol * str_to_sb_symbol(char *cstr)
{
    int x;
    int length;
    sb_symbol *b;

    length = strlen(cstr);
    b = (sb_symbol *) malloc((length+1) * sizeof(sb_symbol));

    for (x = 0; x < length; x++) {
        b[x] = cstr[x];
    }
    b[length] = 0;

    return b;
}

// simply cast a stemmed symbol string back to a
// char for use by C.GoString.
char * sb_symbol_to_char(sb_symbol *b)
{
    return (char *)b;
}

// get the size of the null-terminated list of algorithm names
// there's probably a better way to do this, but my C aint that fly.
int get_num_algorithms() {
    const char ** list;
    const char * entry;
    int count;

    list = sb_stemmer_list();
    count = 0;
    entry = list[0];
    while (entry) {
        count++;
        entry = list[count];
    }

    return count;
}

// get name of the algorithm in position x of the null-terminated
// algorithm list.
const char * get_algorithm_x(int x) {
    const char ** list;

    list = sb_stemmer_list();
    return list[x];
}


*/
import "C"

// get a list of the stemming algorithms supported by this version of libstemmer
func GetSupportedLanguages() []string {
    alg_count := int(C.get_num_algorithms());
    out := make([]string, alg_count)
    for x := 0; x < alg_count; x++ {
        out[x] = C.GoString(C.get_algorithm_x(C.int(x)))
    }
    return out
}


type Stemmer struct {
    stemmer *[0]uint8
    lock *sync.Mutex
}

// internal method for GCing the C allocated stemmer
func (s Stemmer) Close() {
    C.sb_stemmer_delete(s.stemmer)
}

// Create a new stemmer, ready for use with the specified language
func NewStemmer(language string) (*Stemmer, error) {
    clang := C.CString(strings.ToLower(language))
    defer C.free(unsafe.Pointer(clang))
    cchar := C.CString("UTF_8")
    defer C.free(unsafe.Pointer(cchar))
    tmp := C.sb_stemmer_new(clang, cchar)
    if tmp == nil {
        return nil, fmt.Errorf("Unable to create stemmer, please ensure you are using a valid language")
    }
    stemmer := &Stemmer{
        stemmer: tmp,
        lock: new(sync.Mutex),
    }
    return stemmer, nil
}

func (s Stemmer) StemWord(str string) string {
    cstr := C.CString(str)
    defer C.free(unsafe.Pointer(cstr))
    sbs := C.str_to_sb_symbol(cstr)
    defer C.free(unsafe.Pointer(sbs))

    s.lock.Lock()
    stemmed := C.sb_stemmer_stem(s.stemmer, sbs, C.int(len(str)))
    s.lock.Unlock()

    char := C.sb_symbol_to_char(stemmed)
    val := C.GoString(char)

    return val
}

