#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef int (__cdecl *get_info_fn)(int, char *, void *, int);
typedef short (__cdecl *load_fn)(HWND, int, char *, char *);
typedef short (__cdecl *text_to_file_fn)(int, char *, char *, int, int, int, int, int, int, int);
typedef void (__cdecl *unload_fn)(int);

static char *read_text(const char *path) {
    FILE *f = fopen(path, "rb");
    if (!f) return NULL;
    if (fseek(f, 0, SEEK_END) != 0) { fclose(f); return NULL; }
    long size = ftell(f);
    if (size < 0 || size > 1024 * 1024 || fseek(f, 0, SEEK_SET) != 0) {
        fclose(f); return NULL;
    }
    char *text = (char *)calloc((size_t)size + 1, 1);
    if (!text) { fclose(f); return NULL; }
    size_t got = fread(text, 1, (size_t)size, f);
    int failed = ferror(f);
    fclose(f);
    if (failed || got != (size_t)size) { free(text); return NULL; }
    return text;
}

int main(void) {
    HMODULE dll = LoadLibraryA("Z:\\samples\\vt_kat.dll");
    if (!dll) { printf("LOAD_LIBRARY error=%lu\n", (unsigned long)GetLastError()); return 2; }
    FARPROC info_address = GetProcAddress(dll, "VT_GetTTSInfo_ENG");
    FARPROC load_address = GetProcAddress(dll, "VT_LOADTTS_ENG");
    FARPROC file_address = GetProcAddress(dll, "VT_TextToFile_ENG");
    FARPROC unload_address = GetProcAddress(dll, "VT_UNLOADTTS_ENG");
    get_info_fn get_info; load_fn load; text_to_file_fn text_to_file; unload_fn unload;
    memcpy(&get_info, &info_address, sizeof(get_info));
    memcpy(&load, &load_address, sizeof(load));
    memcpy(&text_to_file, &file_address, sizeof(text_to_file));
    memcpy(&unload, &unload_address, sizeof(unload));
    if (!get_info || !load || !text_to_file || !unload) {
        printf("MISSING_EXPORT info=%d load=%d file=%d unload=%d\n",
               get_info != NULL, load != NULL, text_to_file != NULL, unload != NULL);
        FreeLibrary(dll); return 3;
    }
    char date[64] = {0}, dbdir[256] = {0};
    int load_success = -999, speaker = -999, rate = -999;
    printf("INFO build=%d dbdir=%d load_success=%d speaker=%d rate=%d\n",
           get_info(0, NULL, date, sizeof(date)), get_info(3, NULL, dbdir, sizeof(dbdir)),
           get_info(4, NULL, &load_success, sizeof(load_success)),
           get_info(6, NULL, &speaker, sizeof(speaker)),
           get_info(10, NULL, &rate, sizeof(rate)));
    printf("BUILD_DATE value=%s DB_DIRECTORY value=%s load_success_value=%d default_speaker=%d sample_rate=%d\n",
           date, dbdir, load_success, speaker, rate);
    short loaded = load(NULL, -1, NULL, NULL);
    printf("LOAD return=%d expected_success=%d speaker=%d\n", (int)loaded, load_success, speaker);
    if (loaded != load_success) { FreeLibrary(dll); return 4; }
    const char *names[] = {"prose", "numbers", "address"};
    const char *paths[] = {"Z:\\work\\stage17\\input.txt",
                           "Z:\\work\\stage17\\numbers.txt",
                           "Z:\\work\\stage17\\address.txt"};
    int failed = 0;
    for (int i = 0; i < 3; ++i) {
        char *text = read_text(paths[i]);
        if (!text) { printf("FIXTURE name=%s read=failed\n", names[i]); failed = 1; continue; }
        char output[256];
        snprintf(output, sizeof(output), "Z:\\probe\\runtime-original-msi-2006-%s.wav", names[i]);
        short result = text_to_file(4, text, output, speaker, 100, 100, 100, 100, 0, 0);
        printf("SYNTH name=%s return=%d output=%s\n", names[i], (int)result, output);
        if (result != 1) failed = 1;
        free(text);
    }
    unload(speaker);
    FreeLibrary(dll);
    return failed ? 5 : 0;
}
