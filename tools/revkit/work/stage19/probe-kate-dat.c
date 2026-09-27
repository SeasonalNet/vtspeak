/* Sample Kate DAT payloads with the patched standard DLL decoder. */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <wincrypt.h>

enum { HEADER_SIZE = 45, RECORD_SIZE = 19, OUTPUT_CAPACITY = 65536, HASH_SIZE = 32 };
typedef short(__cdecl *load_fn)(HWND, int, char *, char *);
typedef short(__cdecl *synth_fn)(int, char *, char *, int, int, int, int, int, int, int);
typedef void(__cdecl *unload_fn)(int);
typedef unsigned int(__cdecl *decode_fn)(const unsigned char *, unsigned char *, unsigned int, unsigned int *);

static uint16_t u16(const unsigned char *p) { return (uint16_t)p[0] | (uint16_t)((uint16_t)p[1] << 8); }
static uint32_t u32(const unsigned char *p) { return (uint32_t)p[0] | ((uint32_t)p[1] << 8) | ((uint32_t)p[2] << 16) | ((uint32_t)p[3] << 24); }

static int read_file(const char *path, unsigned char **out, size_t *size) {
    FILE *f = fopen(path, "rb");
    if (!f || fseek(f, 0, SEEK_END)) { if (f) fclose(f); return 0; }
    long n = ftell(f);
    if (n <= 0 || fseek(f, 0, SEEK_SET)) { fclose(f); return 0; }
    unsigned char *p = malloc((size_t)n);
    if (!p || fread(p, 1, (size_t)n, f) != (size_t)n) { free(p); fclose(f); return 0; }
    fclose(f); *out = p; *size = (size_t)n; return 1;
}

static int digest(HCRYPTPROV provider, const unsigned char *p, DWORD n, unsigned char out[HASH_SIZE]) {
    HCRYPTHASH h = 0; DWORD size = HASH_SIZE;
    if (!CryptCreateHash(provider, CALG_SHA_256, 0, 0, &h)) return 0;
    int ok = CryptHashData(h, p, n, 0) && CryptGetHashParam(h, HP_HASHVAL, out, &size, 0) && size == HASH_SIZE;
    CryptDestroyHash(h); return ok;
}

int main(int argc, char **argv) {
    int all_units = argc == 2 && !strcmp(argv[1], "--all");
    if (argc > 2 || (argc == 2 && !all_units)) { fprintf(stderr, "usage: probe-kate-dat.exe [--all]\n"); return 2; }
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_kat.dll");
    if (!dll) { fprintf(stderr, "LoadLibrary failed: %lu\n", GetLastError()); return 1; }
    FARPROC p_load = GetProcAddress(dll, "VT_LOADTTS_ENG"), p_synth = GetProcAddress(dll, "VT_TextToFile_ENG"), p_unload = GetProcAddress(dll, "VT_UNLOADTTS_ENG");
    load_fn load; synth_fn synth; unload_fn unload;
    memcpy(&load, &p_load, sizeof(load)); memcpy(&synth, &p_synth, sizeof(synth)); memcpy(&unload, &p_unload, sizeof(unload));
    if (!p_load || !p_synth || !p_unload || load(NULL, -1, NULL, NULL) != 0) { fprintf(stderr, "public DLL initialization failed\n"); return 1; }
    short sr = synth(4, "Hi.", "Z:\\work\\corpus-parity\\kate\\decoder-canary.wav", -1, -1, -1, -1, -1, -1, 0);
    FILE *canary = fopen("Z:\\work\\corpus-parity\\kate\\decoder-canary.wav", "rb");
    unsigned char riff[12] = {0};
    int valid = canary && fread(riff, 1, sizeof(riff), canary) == sizeof(riff) && !memcmp(riff, "RIFF", 4) && !memcmp(riff + 8, "WAVE", 4);
    if (canary) fclose(canary);
    if (sr != 1 || !valid) { fprintf(stderr, "DLL canary failed: result=%d wave=%d\n", sr, valid); unload(-1); return 1; }

    decode_fn decode = (decode_fn)((unsigned char *)dll + 0x1b30);
    HCRYPTPROV provider = 0;
    if (!CryptAcquireContextA(&provider, NULL, NULL, PROV_RSA_AES, CRYPT_VERIFYCONTEXT)) return 1;
    const char *banks[] = {"gen", "gen2", "num", "etc", "alp"};
    unsigned char output[OUTPUT_CAPACITY];
    for (size_t b = 0; b < sizeof(banks) / sizeof(banks[0]); ++b) {
        char path[MAX_PATH]; unsigned char *idx = NULL, *dat = NULL, *upm = NULL; size_t idx_n = 0, dat_n = 0, upm_n = 0;
        snprintf(path, sizeof(path), "Z:\\probe\\data-kate-copy\\M16\\mc_idx_tbl\\unit-%s.idx", banks[b]);
        if (!read_file(path, &idx, &idx_n)) { fprintf(stderr, "cannot read %s index\n", banks[b]); return 1; }
        snprintf(path, sizeof(path), "Z:\\probe\\data-kate-copy\\M16\\dat\\merged-%s.dat", banks[b]);
        if (!read_file(path, &dat, &dat_n)) { fprintf(stderr, "cannot read %s DAT\n", banks[b]); return 1; }
        snprintf(path, sizeof(path), "Z:\\probe\\data-kate-copy\\M16\\dat\\merged-%s.upm", banks[b]);
        if (!read_file(path, &upm, &upm_n)) { fprintf(stderr, "cannot read %s UPM\n", banks[b]); return 1; }
        if (idx_n < HEADER_SIZE || memcmp(idx + 1, "ver.2005\0VoiceText-Eng\0", 23)) { fprintf(stderr, "%s legacy header mismatch\n", banks[b]); return 1; }
        size_t pos = 1u + idx[0] + 2u;
        uint16_t name_n = u16(idx + pos); pos += 2u + name_n + 1u;
        uint32_t count = u32(idx + pos); pos += 6u;
        if (idx_n != pos + (size_t)count * (RECORD_SIZE + 20u)) { fprintf(stderr, "%s index dimensions mismatch\n", banks[b]); return 1; }
        uint32_t sample_units[] = {0, count / 2u, count - 1u};
        uint32_t iterations = all_units ? count : 3u;
        for (uint32_t s = 0; s < iterations; ++s) {
            uint32_t unit = all_units ? s : sample_units[s]; const unsigned char *rec = idx + pos + (size_t)unit * RECORD_SIZE;
            uint32_t off = u32(rec); uint16_t len = u16(rec + 8);
            uint32_t upm_off = u32(rec + 10); uint32_t upm_len = (uint32_t)rec[14] + rec[15] - 1u;
            if ((size_t)off + len > dat_n || (size_t)upm_off + upm_len > upm_n) { fprintf(stderr, "%s/%lu DAT/UPM span invalid\n", banks[b], (unsigned long)unit); return 1; }
            unsigned int expected_bytes = 0;
            for (uint32_t j = 0; j < upm_len; ++j) expected_bytes += 4u * upm[upm_off + j];
            unsigned int out_n = 0;
            unsigned int result = decode(dat + off, output, len, &out_n);
            unsigned char hash[HASH_SIZE];
            if (result != 1 || out_n > sizeof(output) || out_n != expected_bytes || !digest(provider, output, out_n, hash)) { fprintf(stderr, "%s/%lu decoder failed result=%u bytes=%u expected=%u\n", banks[b], (unsigned long)unit, result, out_n, expected_bytes); return 1; }
            printf("%s\t%lu\t%u\t", banks[b], (unsigned long)unit, out_n);
            for (size_t h = 0; h < HASH_SIZE; ++h) printf("%02x", hash[h]);
            putchar('\n');
        }
        free(idx); free(dat); free(upm);
        fprintf(stderr, "%s: decoded %lu/%lu units\n", banks[b], (unsigned long)iterations, (unsigned long)count);
    }
    CryptReleaseContext(provider, 0); unload(-1); FreeLibrary(dll); return 0;
}
