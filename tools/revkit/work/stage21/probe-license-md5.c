#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef void(__cdecl *md5_init_fn)(void *);
typedef void(__cdecl *md5_update_fn)(void *, const void *, uint32_t);
typedef void(__cdecl *md5_final_fn)(unsigned char *, void *);

typedef struct {
    const char *name;
    const char *input;
    uint32_t length;
    const char *expected;
} md5_vector;

static int digest_matches(const unsigned char digest[16], const char *expected) {
    static const char digits[] = "0123456789abcdef";
    char encoded[33];
    for (size_t i = 0; i < 16; ++i) {
        encoded[2 * i] = digits[digest[i] >> 4];
        encoded[2 * i + 1] = digits[digest[i] & 15U];
    }
    encoded[32] = '\0';
    return strcmp(encoded, expected) == 0;
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        printf("HARNESS error=load_library code=%lu\n", (unsigned long)GetLastError());
        return 2;
    }
    uintptr_t base = (uintptr_t)dll;
    md5_init_fn init = (md5_init_fn)(base + 0x142f0U);
    md5_update_fn update = (md5_update_fn)(base + 0x14320U);
    md5_final_fn final = (md5_final_fn)(base + 0x14410U);
    char repeated[128];
    memset(repeated, 'a', sizeof(repeated));
    const md5_vector vectors[] = {
        {"empty", "", 0, "d41d8cd98f00b204e9800998ecf8427e"},
        {"a", "a", 1, "0cc175b9c0f1b6a831c399e269772661"},
        {"abc", "abc", 3, "900150983cd24fb0d6963f7d28e17f72"},
        {"a55", repeated, 55, "ef1772b6dff9a122358552954ad0df65"},
        {"a56", repeated, 56, "3b0c8ac703f828b04c6c197006d17218"},
        {"a63", repeated, 63, "b06521f39153d618550606be297466d5"},
        {"a64", repeated, 64, "014842d480b571495a4a0363793f7367"},
        {"a65", repeated, 65, "c743a45e0d2e6a95cb859adae0248435"},
        {"a127", repeated, 127, "020406e1d05cdc2aa287641f7ae2cc39"},
        {"a128", repeated, 128, "e510683b3f5ffe4093d021808bc6ff70"},
    };
    unsigned char context[128];
    unsigned char digest[16];
    unsigned int failures = 0;
    for (size_t i = 0; i < sizeof(vectors) / sizeof(vectors[0]); ++i) {
        memset(context, 0, sizeof(context));
        memset(digest, 0, sizeof(digest));
        init(context);
        if (vectors[i].length != 0) {
            update(context, vectors[i].input, vectors[i].length);
        }
        final(digest, context);
        int matches = digest_matches(digest, vectors[i].expected);
        failures += !matches;
        printf("MD5_VECTOR name=%s length=%u matches_reference=%d\n", vectors[i].name,
               vectors[i].length, matches);
    }
    printf("MD5_VECTOR_SUMMARY cases=%u failures=%u\n",
           (unsigned int)(sizeof(vectors) / sizeof(vectors[0])), failures);
    FreeLibrary(dll);
    return failures == 0 ? 0 : 1;
}
