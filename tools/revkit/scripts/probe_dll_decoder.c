/* Isolated VoiceText decoder probe; compile as a PE32 program with MinGW. */
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <wincrypt.h>

enum { INDEX_HEADER_SIZE = 45, UNIT_RECORD_SIZE = 19, FEATURE_COLUMNS_SIZE = 21 };
enum { OUTPUT_CAPACITY = 65536, SHA256_BYTES = 32 };

static const char *const banks[] = {"gen", "num", "etc", "alp"};

typedef short(__cdecl *load_tts_fn)(HWND, int, char *, char *);
typedef short(__cdecl *text_to_file_fn)(
    int, char *, char *, int, int, int, int, int, int, int);
typedef void(__cdecl *unload_tts_fn)(int);
typedef unsigned int(__cdecl *decode_fn)(
    const unsigned char *, unsigned char *, unsigned int, unsigned int *);

static int read_file(const char *path, unsigned char **data, size_t *size) {
    FILE *input = fopen(path, "rb");
    if (input == NULL) {
        return 0;
    }
    if (fseek(input, 0, SEEK_END) != 0) {
        fclose(input);
        return 0;
    }
    long length = ftell(input);
    if (length <= 0 || fseek(input, 0, SEEK_SET) != 0) {
        fclose(input);
        return 0;
    }
    unsigned char *buffer = malloc((size_t)length);
    if (buffer == NULL || fread(buffer, 1, (size_t)length, input) != (size_t)length) {
        free(buffer);
        fclose(input);
        return 0;
    }
    fclose(input);
    *data = buffer;
    *size = (size_t)length;
    return 1;
}

static uint16_t read_u16(const unsigned char *source) {
    return (uint16_t)source[0] | (uint16_t)((uint16_t)source[1] << 8);
}

static uint32_t read_u32(const unsigned char *source) {
    return (uint32_t)source[0] | ((uint32_t)source[1] << 8) |
           ((uint32_t)source[2] << 16) | ((uint32_t)source[3] << 24);
}

static int digest_pcm(HCRYPTPROV provider, const unsigned char *pcm, unsigned int size,
                      unsigned char digest[SHA256_BYTES]) {
    HCRYPTHASH hash = 0;
    DWORD digest_size = SHA256_BYTES;
    if (!CryptCreateHash(provider, CALG_SHA_256, 0, 0, &hash)) {
        return 0;
    }
    int ok = CryptHashData(hash, pcm, size, 0) &&
             CryptGetHashParam(hash, HP_HASHVAL, digest, &digest_size, 0) &&
             digest_size == SHA256_BYTES;
    CryptDestroyHash(hash);
    return ok;
}

static int scan_bank(const char *bank, decode_fn decode, HCRYPTPROV provider,
                     FILE *result, unsigned char output[OUTPUT_CAPACITY], uint32_t *count_out) {
    char path[MAX_PATH];
    unsigned char *index = NULL;
    unsigned char *dat = NULL;
    unsigned char *upm = NULL;
    size_t index_size = 0, dat_size = 0, upm_size = 0;
    snprintf(path, sizeof(path), "Z:\\voice-data\\paul\\M16\\mc_idx_tbl\\unit-%s.idx", bank);
    if (!read_file(path, &index, &index_size)) {
        fprintf(stderr, "%s: cannot read index\n", bank);
        goto failure;
    }
    snprintf(path, sizeof(path), "Z:\\voice-data\\paul\\M16\\dat\\merged-%s.dat", bank);
    if (!read_file(path, &dat, &dat_size)) {
        fprintf(stderr, "%s: cannot read DAT bank\n", bank);
        goto failure;
    }
    snprintf(path, sizeof(path), "Z:\\voice-data\\paul\\M16\\dat\\merged-%s.upm", bank);
    if (!read_file(path, &upm, &upm_size)) {
        fprintf(stderr, "%s: cannot read UPM bank\n", bank);
        goto failure;
    }
    if (index_size < INDEX_HEADER_SIZE ||
        memcmp(index + 1, "ver.2013\0VoiceText-Eng\0", 23) != 0) {
        fprintf(stderr, "%s: unsupported index header\n", bank);
        goto failure;
    }
    uint32_t count = read_u32(index + 39);
    uint16_t stride = read_u16(index + 43);
    if (stride != UNIT_RECORD_SIZE ||
        index_size != INDEX_HEADER_SIZE + (size_t)count *
                          (UNIT_RECORD_SIZE + FEATURE_COLUMNS_SIZE)) {
        fprintf(stderr, "%s: unsupported index dimensions\n", bank);
        goto failure;
    }

    for (uint32_t unit = 0; unit < count; ++unit) {
        const unsigned char *record = index + INDEX_HEADER_SIZE + (size_t)unit * stride;
        uint32_t dat_offset = read_u32(record);
        uint16_t dat_length = read_u16(record + 8);
        uint32_t upm_offset = read_u32(record + 10);
        unsigned int upm_length = (unsigned int)record[14] + record[15] - 1;
        if (dat_offset > dat_size || dat_length > dat_size - dat_offset ||
            upm_offset > upm_size || upm_length > upm_size - upm_offset) {
            fprintf(stderr, "%s unit %lu: payload span out of bounds\n",
                    bank, (unsigned long)unit);
            goto failure;
        }
        unsigned int expected_bytes = 0;
        for (unsigned int period = 0; period < upm_length; ++period) {
            expected_bytes += 4u * upm[upm_offset + period];
        }
        if (expected_bytes == 0 || expected_bytes > OUTPUT_CAPACITY) {
            fprintf(stderr, "%s unit %lu: invalid expected PCM size\n",
                    bank, (unsigned long)unit);
            goto failure;
        }
        unsigned int output_bytes = 0;
        unsigned int decode_result =
            decode(dat + dat_offset, output, dat_length, &output_bytes);
        if (decode_result != 1 || output_bytes != expected_bytes) {
            fprintf(stderr, "%s unit %lu: decoder result=%u bytes=%u expected=%u\n",
                    bank, (unsigned long)unit, decode_result, output_bytes, expected_bytes);
            goto failure;
        }
        unsigned char digest[SHA256_BYTES];
        if (!digest_pcm(provider, output, output_bytes, digest) ||
            fwrite(&output_bytes, sizeof(output_bytes), 1, result) != 1 ||
            fwrite(digest, sizeof(digest), 1, result) != 1) {
            fprintf(stderr, "%s unit %lu: digest or result write failed\n",
                    bank, (unsigned long)unit);
            goto failure;
        }
        if (unit != 0 && unit % 100000 == 0) {
            printf("%s: decoded %lu/%lu units\n", bank,
                   (unsigned long)unit, (unsigned long)count);
            fflush(stdout);
        }
    }
    *count_out = count;
    free(index);
    free(dat);
    free(upm);
    return 1;

failure:
    free(index);
    free(dat);
    free(upm);
    return 0;
}

static int run_corpus(decode_fn decode) {
    HCRYPTPROV provider = 0;
    if (!CryptAcquireContextA(&provider, NULL, NULL, PROV_RSA_AES, CRYPT_VERIFYCONTEXT)) {
        fprintf(stderr, "SHA-256 provider unavailable: %lu\n", GetLastError());
        return 0;
    }
    FILE *result = fopen("Z:\\work\\corpus-parity\\dll-sha256.bin", "wb");
    if (result == NULL || fwrite("VTDIGEST1", 1, 9, result) != 9) {
        fprintf(stderr, "cannot open corpus result\n");
        if (result != NULL) {
            fclose(result);
        }
        CryptReleaseContext(provider, 0);
        return 0;
    }
    unsigned char *output = malloc(OUTPUT_CAPACITY);
    if (output == NULL) {
        fprintf(stderr, "cannot allocate decoder output buffer\n");
        fclose(result);
        CryptReleaseContext(provider, 0);
        return 0;
    }
    int ok = 1;
    for (size_t bank = 0; bank < sizeof(banks) / sizeof(banks[0]); ++bank) {
        uint32_t count = 0;
        if (!scan_bank(banks[bank], decode, provider, result, output, &count)) {
            ok = 0;
            break;
        }
        printf("%s: complete %lu units\n", banks[bank], (unsigned long)count);
        fflush(stdout);
    }
    free(output);
    if (fclose(result) != 0) {
        ok = 0;
    }
    CryptReleaseContext(provider, 0);
    return ok;
}

int main(int argc, char **argv) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_pau.dll");
    if (dll == NULL) {
        fprintf(stderr, "DLL load failed: %lu\n", GetLastError());
        return 1;
    }
    _Static_assert(sizeof(load_tts_fn) == sizeof(FARPROC), "function pointer width");
    _Static_assert(sizeof(text_to_file_fn) == sizeof(FARPROC), "function pointer width");
    _Static_assert(sizeof(unload_tts_fn) == sizeof(FARPROC), "function pointer width");
    FARPROC load_export = GetProcAddress(dll, "VT_LOADTTS_ENG");
    FARPROC synthesis_export = GetProcAddress(dll, "VT_TextToFile_ENG");
    FARPROC unload_export = GetProcAddress(dll, "VT_UNLOADTTS_ENG");
    load_tts_fn load_tts;
    text_to_file_fn text_to_file;
    unload_tts_fn unload_tts;
    memcpy(&load_tts, &load_export, sizeof(load_tts));
    memcpy(&text_to_file, &synthesis_export, sizeof(text_to_file));
    memcpy(&unload_tts, &unload_export, sizeof(unload_tts));
    if (load_tts == NULL || text_to_file == NULL || unload_tts == NULL) {
        fprintf(stderr, "required public export missing\n");
        return 1;
    }

    short load_result = load_tts(NULL, -1, NULL, NULL);
    printf("public load result=%d\n", load_result);
    short synth_result = text_to_file(
        4, "Hi.", "Z:\\work\\corpus-parity\\canary.wav", -1, -1, -1, -1, -1, -1, 0);
    FILE *canary = fopen("Z:\\work\\corpus-parity\\canary.wav", "rb");
    unsigned char header[12] = {0};
    int canary_valid = canary != NULL && fread(header, 1, sizeof(header), canary) == sizeof(header) &&
                       memcmp(header, "RIFF", 4) == 0 && memcmp(header + 8, "WAVE", 4) == 0;
    if (canary != NULL) {
        fclose(canary);
    }
    printf("public synthesis result=%d valid_wave=%d\n", synth_result, canary_valid);
    if (load_result != 0 || synth_result != 1 || !canary_valid) {
        unload_tts(-1);
        return 1;
    }

    unsigned char *payload = NULL;
    size_t payload_size = 0;
    if (!read_file("Z:\\work\\stage2-copy\\candidate-00.dat", &payload, &payload_size)) {
        fprintf(stderr, "captured DAT payload unavailable\n");
        unload_tts(-1);
        return 1;
    }
    unsigned char output[65536] = {0};
    unsigned int output_bytes = 0;
    decode_fn decode = (decode_fn)((unsigned char *)dll + 0x1b30);
    unsigned int decode_result = decode(payload, output, (unsigned int)payload_size, &output_bytes);
    free(payload);
    printf("decoder result=%u bytes=%u\n", decode_result, output_bytes);
    if (decode_result != 1 || output_bytes > sizeof(output)) {
        unload_tts(-1);
        return 1;
    }
    FILE *pcm = fopen("Z:\\work\\corpus-parity\\smoke.pcm", "wb");
    if (pcm == NULL || fwrite(output, 1, output_bytes, pcm) != output_bytes) {
        fprintf(stderr, "could not write probe PCM\n");
        if (pcm != NULL) {
            fclose(pcm);
        }
        unload_tts(-1);
        return 1;
    }
    fclose(pcm);
    int corpus_ok = 1;
    if (argc == 2 && strcmp(argv[1], "--corpus") == 0) {
        corpus_ok = run_corpus(decode);
    } else if (argc != 1) {
        fprintf(stderr, "usage: probe_dll_decoder.exe [--corpus]\n");
        corpus_ok = 0;
    }
    unload_tts(-1);
    FreeLibrary(dll);
    return corpus_ok ? 0 : 1;
}
